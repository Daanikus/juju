// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cloud

// Generated code - do not edit.

const fallbackPublicCloudInfo = `# DO NOT EDIT, will be overwritten, use "juju update-clouds" to refresh.
clouds:
  aws:
    type: ec2
    description: Amazon Web Services
    auth-types: [ access-key ]
    regions:
      us-east-1:
        endpoint: https://ec2.us-east-1.amazonaws.com
      us-east-2:
        endpoint: https://ec2.us-east-2.amazonaws.com
      us-west-1:
        endpoint: https://ec2.us-west-1.amazonaws.com
      us-west-2:
        endpoint: https://ec2.us-west-2.amazonaws.com
      ca-central-1:
        endpoint: https://ec2.ca-central-1.amazonaws.com
      eu-west-1:
        endpoint: https://ec2.eu-west-1.amazonaws.com
      eu-west-2:
        endpoint: https://ec2.eu-west-2.amazonaws.com
      eu-west-3:
        endpoint: https://ec2.eu-west-3.amazonaws.com
      eu-central-1:
        endpoint: https://ec2.eu-central-1.amazonaws.com
      ap-south-1:
        endpoint: https://ec2.ap-south-1.amazonaws.com
      ap-southeast-1:
        endpoint: https://ec2.ap-southeast-1.amazonaws.com
      ap-southeast-2:
        endpoint: https://ec2.ap-southeast-2.amazonaws.com
      ap-northeast-1:
        endpoint: https://ec2.ap-northeast-1.amazonaws.com
      ap-northeast-2:
        endpoint: https://ec2.ap-northeast-2.amazonaws.com
      sa-east-1:
        endpoint: https://ec2.sa-east-1.amazonaws.com
  aws-china:
    type: ec2
    description: Amazon China
    auth-types: [ access-key ]
    regions:
      cn-north-1:
        endpoint: https://ec2.cn-north-1.amazonaws.com.cn
  aws-gov:
    type: ec2
    description: Amazon (USA Government)
    auth-types: [ access-key ]
    regions:
      us-gov-west-1:
        endpoint: https://ec2.us-gov-west-1.amazonaws.com
  google:
    type: gce
    description: Google Cloud Platform
    auth-types: [ jsonfile, oauth2 ]
    regions:
      us-east1:
        endpoint: https://www.googleapis.com
      us-east4:
        endpoint: https://www.googleapis.com
      us-central1:
        endpoint: https://www.googleapis.com
      us-west1:
        endpoint: https://www.googleapis.com
      europe-west1:
        endpoint: https://www.googleapis.com
      europe-west2:
        endpoint: https://www.googleapis.com
      europe-west3:
        endpoint: https://www.googleapis.com
      asia-east1:
        endpoint: https://www.googleapis.com
      asia-northeast1:
        endpoint: https://www.googleapis.com
      asia-southeast1:
        endpoint: https://www.googleapis.com
      asia-south1:
        endpoint: https://www.googleapis.com
      australia-southeast1:
        endpoint: https://www.googleapis.com
      southamerica-east1:
        endpoint: https://www.googleapis.com
  azure:
    type: azure
    description: Microsoft Azure
    auth-types: [ interactive, service-principal-secret ]
    regions:
      centralus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      eastus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      eastus2:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      northcentralus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      southcentralus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      westcentralus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      westus:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      westus2:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      northeurope:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      westeurope:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      eastasia:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      southeastasia:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      japaneast:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      japanwest:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      brazilsouth:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      australiaeast:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      australiasoutheast:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      centralindia:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      southindia:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      westindia:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      canadacentral:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      canadaeast:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      uksouth:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      ukwest:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net
      koreacentral:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net        
      koreasouth:
        endpoint: https://management.azure.com
        storage-endpoint: https://core.windows.net
        identity-endpoint: https://graph.windows.net        
  azure-china:
    type: azure
    description: Microsoft Azure China
    auth-types: [ interactive, service-principal-secret ]
    regions:
      chinaeast:
        endpoint: https://management.chinacloudapi.cn
        storage-endpoint: https://core.chinacloudapi.cn
        identity-endpoint: https://graph.chinacloudapi.cn
      chinanorth:
        endpoint: https://management.chinacloudapi.cn
        storage-endpoint: https://core.chinacloudapi.cn
        identity-endpoint: https://graph.chinacloudapi.cn
  rackspace:
    type: rackspace
    description: Rackspace Cloud
    auth-types: [ userpass ]
    endpoint: https://identity.api.rackspacecloud.com/v2.0
    regions:
      dfw:
        endpoint: https://identity.api.rackspacecloud.com/v2.0
      ord:
        endpoint: https://identity.api.rackspacecloud.com/v2.0
      iad:
        endpoint: https://identity.api.rackspacecloud.com/v2.0
      lon:
        endpoint: https://lon.identity.api.rackspacecloud.com/v2.0
      syd:
        endpoint: https://identity.api.rackspacecloud.com/v2.0
      hkg:
        endpoint: https://identity.api.rackspacecloud.com/v2.0
  joyent:
    type: joyent
    description: Joyent Cloud
    auth-types: [ userpass ]
    regions:
      eu-ams-1:
        endpoint: https://eu-ams-1.api.joyentcloud.com
      us-sw-1:
        endpoint: https://us-sw-1.api.joyentcloud.com
      us-east-1:
        endpoint: https://us-east-1.api.joyentcloud.com
      us-east-2:
        endpoint: https://us-east-2.api.joyentcloud.com
      us-east-3:
        endpoint: https://us-east-3.api.joyentcloud.com
      us-west-1:
        endpoint: https://us-west-1.api.joyentcloud.com
  cloudsigma:
    type: cloudsigma
    description: CloudSigma Cloud
    auth-types: [ userpass ]
    regions:
      hnl:
        endpoint: https://hnl.cloudsigma.com/api/2.0/
      mia:
        endpoint: https://mia.cloudsigma.com/api/2.0/
      sjc:
        endpoint: https://sjc.cloudsigma.com/api/2.0/
      wdc:
        endpoint: https://wdc.cloudsigma.com/api/2.0/
      zrh:
        endpoint: https://zrh.cloudsigma.com/api/2.0/
  oracle:
    type: oci
    description: Oracle Cloud Infrastructure
    auth-types: [ httpsig ]
    regions:
      us-phoenix-1:
        endpoint: https://iaas.us-phoenix-1.oraclecloud.com
      us-ashburn-1:
        endpoint: https://iaas.us-ashburn-1.oraclecloud.com
      eu-frankfurt-1:
        endpoint: https://iaas.eu-frankfurt-1.oraclecloud.com
      uk-london-1:
        endpoint: https://iaas.uk-london-1.oraclecloud.com
`
