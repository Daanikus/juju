// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package remotestate

import (
	"sync"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names"
	"launchpad.net/tomb"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state/watcher"
	"github.com/juju/juju/worker"
	"github.com/juju/juju/worker/leadership"
)

var logger = loggo.GetLogger("juju.worker.uniter.remotestate")

// RemoteStateWatcher collects unit, service, and service config information
// from separate state watchers, and updates a Snapshot which is sent on a
// channel upon change.
type RemoteStateWatcher struct {
	st                         State
	unit                       Unit
	service                    Service
	relations                  map[names.RelationTag]*relationUnitsWatcher
	relationUnitsChanges       chan relationUnitsChange
	storageAttachementWatchers map[names.StorageTag]*storageAttachmentWatcher
	storageAttachment          chan StorageSnapshotEvent
	leadershipTracker          leadership.Tracker
	tomb                       tomb.Tomb

	out     chan struct{}
	mu      sync.Mutex
	current Snapshot
}

// NewWatcher returns a RemoteStateWatcher that handles state changes pertaining to the
// supplied unit.
func NewWatcher(st State, leadershipTracker leadership.Tracker, unitTag names.UnitTag) (*RemoteStateWatcher, error) {
	w := &RemoteStateWatcher{
		st:                         st,
		relations:                  make(map[names.RelationTag]*relationUnitsWatcher),
		relationUnitsChanges:       make(chan relationUnitsChange),
		storageAttachementWatchers: make(map[names.StorageTag]*storageAttachmentWatcher),
		storageAttachment:          make(chan StorageSnapshotEvent),
		leadershipTracker:          leadershipTracker,
		out:                        make(chan struct{}),
		current: Snapshot{
			Relations: make(map[int]RelationSnapshot),
			Storage:   make(map[names.StorageTag]StorageSnapshot),
		},
	}
	if err := w.init(unitTag); err != nil {
		return nil, errors.Trace(err)
	}
	go func() {
		defer w.tomb.Done()
		err := w.loop(unitTag)
		logger.Errorf("remote state watcher exited: %v", err)
		w.tomb.Kill(err)
	}()
	return w, nil
}

func (w *RemoteStateWatcher) Stop() error {
	w.tomb.Kill(nil)
	return w.tomb.Wait()
}

func (w *RemoteStateWatcher) Dead() <-chan struct{} {
	return w.tomb.Dead()
}

func (w *RemoteStateWatcher) Wait() error {
	return w.tomb.Wait()
}

func (w *RemoteStateWatcher) Kill() {
	w.tomb.Kill(nil)
}

func (w *RemoteStateWatcher) RemoteStateChanged() <-chan struct{} {
	return w.out
}

func (w *RemoteStateWatcher) Snapshot() Snapshot {
	w.mu.Lock()
	defer w.mu.Unlock()
	snapshot := w.current
	snapshot.Relations = make(map[int]RelationSnapshot)
	for id, relationSnapshot := range w.current.Relations {
		snapshot.Relations[id] = relationSnapshot
	}
	snapshot.Storage = make(map[names.StorageTag]StorageSnapshot)
	for tag, storageSnapshot := range w.current.Storage {
		snapshot.Storage[tag] = storageSnapshot
	}
	return snapshot
}

func (w *RemoteStateWatcher) ClearResolvedMode() {
	w.mu.Lock()
	w.current.ResolvedMode = params.ResolvedNone
	w.mu.Unlock()
}

func (w *RemoteStateWatcher) init(unitTag names.UnitTag) (err error) {
	// TODO(dfc) named return value is a time bomb
	// TODO(axw) move this logic.
	defer func() {
		if params.IsCodeNotFoundOrCodeUnauthorized(err) {
			err = worker.ErrTerminateAgent
		}
	}()
	if w.unit, err = w.st.Unit(unitTag); err != nil {
		return err
	}
	w.service, err = w.unit.Service()
	if err != nil {
		return err
	}
	return nil
}

func (w *RemoteStateWatcher) loop(unitTag names.UnitTag) (err error) {
	var requiredEvents int

	var seenUnitChange bool
	unitw, err := w.unit.Watch()
	if err != nil {
		return err
	}
	defer watcher.Stop(unitw, &w.tomb)
	requiredEvents++

	var seenServiceChange bool
	servicew, err := w.service.Watch()
	if err != nil {
		return err
	}
	defer watcher.Stop(servicew, &w.tomb)
	requiredEvents++

	var seenConfigChange bool
	configw, err := w.unit.WatchConfigSettings()
	if err != nil {
		return err
	}
	defer watcher.Stop(configw, &w.tomb)
	requiredEvents++

	var seenRelationsChange bool
	relationsw, err := w.service.WatchRelations()
	if err != nil {
		return err
	}
	defer watcher.Stop(relationsw, &w.tomb)
	requiredEvents++

	var seenAddressesChange bool
	addressesw, err := w.unit.WatchAddresses()
	if err != nil {
		return err
	}
	defer watcher.Stop(addressesw, &w.tomb)
	requiredEvents++

	var seenStorageChange bool
	storagew, err := w.unit.WatchStorage()
	if err != nil {
		return err
	}
	defer watcher.Stop(storagew, &w.tomb)
	requiredEvents++

	var seenLeaderSettingsChange bool
	leaderSettingsw, err := w.service.WatchLeadershipSettings()
	if err != nil {
		return err
	}
	defer watcher.Stop(leaderSettingsw, &w.tomb)
	requiredEvents++

	var seenLeadershipChange bool
	// There's no watcher for this per se; we wait on a channel
	// returned by the leadership tracker.
	requiredEvents++

	var eventsObserved int
	observedEvent := func(flag *bool) {
		if !*flag {
			*flag = true
			eventsObserved++
		}
	}

	// fire will, once the first event for each watcher has
	// been observed, send a signal on the out channel.
	fire := func() {
		if eventsObserved != requiredEvents {
			return
		}
		select {
		case w.out <- struct{}{}:
		default:
		}
	}

	defer func() {
		for _, ruw := range w.relations {
			watcher.Stop(ruw, &w.tomb)
		}
	}()

	// Check the initial leadership status, and then we can flip-flop
	// waiting on leader or minion to trigger the changed event.
	var waitLeader, waitMinion <-chan struct{}
	claimLeader := w.leadershipTracker.ClaimLeader()
	select {
	case <-w.tomb.Dying():
		return tomb.ErrDying
	case <-claimLeader.Ready():
		isLeader := claimLeader.Wait()
		w.leadershipChanged(isLeader)
		if isLeader {
			waitMinion = w.leadershipTracker.WaitMinion().Ready()
		} else {
			waitLeader = w.leadershipTracker.WaitLeader().Ready()
		}
		observedEvent(&seenLeadershipChange)
	}

	for {
		select {
		case <-w.tomb.Dying():
			return tomb.ErrDying

		case _, ok := <-unitw.Changes():
			logger.Debugf("got unit change")
			if !ok {
				return watcher.EnsureErr(unitw)
			}
			if err := w.unitChanged(); err != nil {
				return err
			}
			observedEvent(&seenUnitChange)

		case _, ok := <-servicew.Changes():
			logger.Debugf("got service change")
			if !ok {
				return watcher.EnsureErr(servicew)
			}
			if err := w.serviceChanged(); err != nil {
				return err
			}
			observedEvent(&seenServiceChange)

		case _, ok := <-configw.Changes():
			logger.Debugf("got config change")
			if !ok {
				return watcher.EnsureErr(configw)
			}
			if err := w.configChanged(); err != nil {
				return err
			}
			observedEvent(&seenConfigChange)

		case _, ok := <-addressesw.Changes():
			logger.Debugf("got address change")
			if !ok {
				return watcher.EnsureErr(addressesw)
			}
			if err := w.addressesChanged(); err != nil {
				return err
			}
			observedEvent(&seenAddressesChange)

		case _, ok := <-leaderSettingsw.Changes():
			logger.Debugf("got leader settings change: ok=%t", ok)
			if !ok {
				return watcher.EnsureErr(leaderSettingsw)
			}
			if err := w.leaderSettingsChanged(); err != nil {
				return err
			}
			observedEvent(&seenLeaderSettingsChange)

		case keys, ok := <-relationsw.Changes():
			logger.Debugf("got relations change")
			if !ok {
				return watcher.EnsureErr(relationsw)
			}
			if err := w.relationsChanged(keys); err != nil {
				return err
			}
			observedEvent(&seenRelationsChange)

		case keys, ok := <-storagew.Changes():
			logger.Debugf("got storage change: %v", keys)
			if !ok {
				return watcher.EnsureErr(storagew)
			}
			if err := w.storageChanged(keys); err != nil {
				return err
			}
			observedEvent(&seenStorageChange)

		case <-waitMinion:
			logger.Debugf("got leadership change: minion")
			if err := w.leadershipChanged(false); err != nil {
				return err
			}
			waitMinion = nil
			waitLeader = w.leadershipTracker.WaitLeader().Ready()

		case <-waitLeader:
			logger.Debugf("got leadership change: leader")
			if err := w.leadershipChanged(true); err != nil {
				return err
			}
			waitLeader = nil
			waitMinion = w.leadershipTracker.WaitMinion().Ready()

		case event, ok := <-w.storageAttachment:
			logger.Debugf("storage %v snapshot event %v", event.Tag, event)
			if ok {
				if err := w.storageAttachmentChanged(event); err != nil {
					return err
				}
			}

		case change := <-w.relationUnitsChanges:
			logger.Debugf("got a relation units change")
			if err := w.relationUnitsChanged(change); err != nil {
				return err
			}
		}

		// Something changed.
		fire()
	}
}

// unitChanged responds to changes in the unit.
func (w *RemoteStateWatcher) unitChanged() error {
	if err := w.unit.Refresh(); err != nil {
		return err
	}
	resolved, err := w.unit.Resolved()
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.current.Life = w.unit.Life()
	w.current.ResolvedMode = resolved
	return nil
}

// serviceChanged responds to changes in the service.
func (w *RemoteStateWatcher) serviceChanged() error {
	if err := w.service.Refresh(); err != nil {
		return err
	}
	url, force, err := w.service.CharmURL()
	if err != nil {
		return err
	}
	w.mu.Lock()
	w.current.CharmURL = url
	w.current.ForceCharmUpgrade = force
	w.mu.Unlock()
	return nil
}

func (w *RemoteStateWatcher) configChanged() error {
	w.mu.Lock()
	w.current.ConfigVersion++
	w.mu.Unlock()
	return nil
}

func (w *RemoteStateWatcher) addressesChanged() error {
	w.mu.Lock()
	w.current.ConfigVersion++
	w.mu.Unlock()
	return nil
}

func (w *RemoteStateWatcher) leaderSettingsChanged() error {
	w.mu.Lock()
	w.current.LeaderSettingsVersion++
	w.mu.Unlock()
	return nil
}

func (w *RemoteStateWatcher) leadershipChanged(isLeader bool) error {
	w.mu.Lock()
	w.current.Leader = isLeader
	w.mu.Unlock()
	return nil
}

// relationsChanged responds to service relation changes.
func (w *RemoteStateWatcher) relationsChanged(keys []string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, key := range keys {
		relationTag := names.NewRelationTag(key)
		rel, err := w.st.Relation(relationTag)
		if params.IsCodeNotFoundOrCodeUnauthorized(err) {
			// If it's actually gone, this unit cannot have entered
			// scope, and therefore never needs to know about it.
			if ruw, ok := w.relations[relationTag]; ok {
				if err := ruw.Stop(); err != nil {
					return errors.Trace(err)
				}
				delete(w.relations, relationTag)
				delete(w.current.Relations, ruw.relationId)
			}
		} else if err != nil {
			return err
		} else {
			if _, ok := w.relations[relationTag]; ok {
				relationSnapshot := w.current.Relations[rel.Id()]
				relationSnapshot.Life = rel.Life()
				w.current.Relations[rel.Id()] = relationSnapshot
				continue
			}
			in, err := w.st.WatchRelationUnits(relationTag, w.unit.Tag())
			if err != nil {
				return errors.Trace(err)
			}
			relationSnapshot := RelationSnapshot{
				Life:    rel.Life(),
				Members: make(map[string]int64),
			}
			select {
			case <-w.tomb.Dying():
				return tomb.ErrDying
			case change, ok := <-in.Changes():
				if !ok {
					return watcher.EnsureErr(in)
				}
				for unit, settings := range change.Changed {
					relationSnapshot.Members[unit] = settings.Version
				}
			}
			w.current.Relations[rel.Id()] = relationSnapshot
			w.relations[relationTag] = newRelationUnitsWatcher(
				rel.Id(), in, w.relationUnitsChanges,
			)
		}
	}
	return nil
}

// relationUnitsChanged responds to relation units changes.
func (w *RemoteStateWatcher) relationUnitsChanged(change relationUnitsChange) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	snapshot, ok := w.current.Relations[change.relationId]
	if !ok {
		return nil
	}
	for unit, settings := range change.Changed {
		snapshot.Members[unit] = settings.Version
	}
	for _, unit := range change.Departed {
		delete(snapshot.Members, unit)
	}
	return nil
}

// storageAttachmentChanged responds to storage attachment changes.
func (w *RemoteStateWatcher) storageAttachmentChanged(event StorageSnapshotEvent) error {
	w.mu.Lock()
	if event.remove {
		delete(w.current.Storage, event.StorageSnapshot.Tag)
	} else {
		w.current.Storage[event.StorageSnapshot.Tag] = event.StorageSnapshot
	}
	w.mu.Unlock()
	return nil
}

// storageChanged responds to unit storage changes.
func (w *RemoteStateWatcher) storageChanged(keys []string) error {
	tags := make([]names.StorageTag, len(keys))
	for i, key := range keys {
		tags[i] = names.NewStorageTag(key)
	}
	ids := make([]params.StorageAttachmentId, len(keys))
	for i, tag := range tags {
		ids[i] = params.StorageAttachmentId{
			StorageTag: tag.String(),
			UnitTag:    w.unit.Tag().String(),
		}
	}
	results, err := w.st.StorageAttachmentLife(ids)
	if err != nil {
		return errors.Trace(err)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, result := range results {
		logger.Debugf("storage result %v", result)
		tag := tags[i]
		// If the storage is alive and present, we start a watcher for it.
		if result.Error == nil {
			storageSnapshot, ok := w.current.Storage[tag]
			if !ok {
				// TODO(wallyworld) - this should be there after consuming initial watcher event
				storageSnapshot = StorageSnapshot{
					Tag:  tag,
					Life: result.Life,
				}
			}
			storageSnapshot.Life = result.Life
			w.current.Storage[tag] = storageSnapshot

			if err := w.startStorageAttachmentWatcher(tag); err != nil {
				return errors.Annotatef(
					err, "starting watcher of %s attachment",
					names.ReadableString(tag),
				)
			}
		} else if params.IsCodeNotFound(result.Error) {
			delete(w.current.Storage, tag)
			if err := w.stopStorageAttachmentWatcher(tag); err != nil {
				return errors.Annotatef(
					err, "stopping watcher of %s attachment",
					names.ReadableString(tag),
				)
			}
		} else {
			logger.Errorf("error getting life of %v attachment: %v", names.ReadableString(tag), err)
			return errors.Annotatef(
				result.Error, "getting life of %s attachment",
				names.ReadableString(tag),
			)
		}
	}
	logger.Debugf("storage change. snapshot is %v", w.current.Storage)
	return nil
}

func (w *RemoteStateWatcher) startStorageAttachmentWatcher(tag names.StorageTag) error {
	if _, ok := w.storageAttachementWatchers[tag]; ok {
		return nil
	}
	logger.Debugf("starting storage attachment watcher for %v", tag)
	watcher, err := newStorageAttachmentWatcher(w.st, w.unit.Tag(), tag, w.storageAttachment)
	if err != nil {
		return errors.Trace(err)
	}
	w.storageAttachementWatchers[tag] = watcher
	return nil
}

func (w *RemoteStateWatcher) stopStorageAttachmentWatcher(tag names.StorageTag) error {
	if watcher, ok := w.storageAttachementWatchers[tag]; !ok {
		return nil
	} else {
		delete(w.storageAttachementWatchers, tag)
		logger.Debugf("stopping storage attachment watcher for %v", tag)
		return watcher.Stop()
	}
	return nil
}
