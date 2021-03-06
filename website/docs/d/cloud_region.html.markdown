---
layout: "ovh"
page_title: "OVH: cloud_region"
sidebar_current: "docs-ovh-datasource-cloud-region"
description: |-
  Get information & status of a region associated with a public cloud project.
---

# cloud_region

Use this data source to retrieve information about a region associated with a
public cloud project. The region must be associated with the project.

## Example Usage

```hcl
data "ovh_cloud_region" "GRA1" {
   project_id = "XXXXXX"
   name = "GRA1"
}
```

## Argument Reference


* `project_id` - (Required) The id of the public cloud project. If omitted,
    the `OVH_PROJECT_ID` environment variable is used.

* `name` - (Required) The name of the region associated with the public cloud
project.

## Attributes Reference

`id` is set to the ID of the project concatenated with the name of the region.
In addition, the following attributes are exported:

* `continent_code` - the code of the geographic continent the region is running.
E.g.: EU for Europe, US for America...
* `datacenter_location` - The location code of the datacenter.
E.g.: "GRA", meaning Gravelines, for region "GRA1"
* `continentCode` - (Deprecated) Use `continent_code` instead.
* `datacenterLocation` - (Deprecated) Use `datacenter_location` instead.
* `services` - The list of public cloud services running within the region
  * `name` - the name of the public cloud service
  * `status` - the status of the service
