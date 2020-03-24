---
layout: "time"
page_title: "Time: time_rotating"
description: |-
  Manages a rotating time resource.
---

# Resource: time_rotating

Manages a rotating time resource, which keeps a rotating UTC timestamp stored in the Terraform state and proposes resource recreation when the locally sourced current time is beyond the rotation time. This rotation only occurs when Terraform is executed, meaning there will be drift between the rotation timestamp and actual rotation. The new rotation timestamp offset includes this drift. This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html) by only forcing a new value on the set cadence.

-> Further manipulation of incoming or outgoing values can be accomplished with the [`formatdate()` function](https://www.terraform.io/docs/configuration/functions/formatdate.html) and the [`timeadd()` function](https://www.terraform.io/docs/configuration/functions/timeadd.html).

## Example Usage

This example configuration will rotate (destroy/create) the resource every 30 days.

```hcl
resource "time_rotating" "example" {
  rotation_days = 30
}
```

## Argument Reference

~> **NOTE:** At least one of the `rotation_` arguments must be configured.

The following arguments are optional:

* `base_rfc3339` - (Optional) Configure the base timestamp with an UTC [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.
* `keepers` - (Optional) Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. These conditions recreate the resource in addition to other rotation arguments. See [the main provider documentation](../index.html) for more information.
* `rotation_days` - (Optional) Number of days to add to the base timestamp to configure the rotation timestamp. When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.
* `rotation_hours` - (Optional) Number of hours to add to the base timestamp to configure the rotation timestamp. When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.
* `rotation_minutes` - (Optional) Number of minutes to add to the base timestamp to configure the rotation timestamp. When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.
* `rotation_months` - (Optional) Number of months to add to the base timestamp to configure the rotation timestamp. When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.
* `rotation_rfc3339` - (Optional) Configure the rotation timestamp with an UTC [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`). When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.
* `rotation_years` - (Optional) Number of years to add to the base timestamp to configure the rotation timestamp. When the current time has passed the rotation timestamp, the resource will trigger recreation. Conflicts with other `rotation_` arguments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `day` - Number day of timestamp.
* `hour` - Number hour of timestamp.
* `id` - UTC RFC3339 format of the base timestamp, e.g. `2020-02-12T06:36:13Z`.
* `minute` - Number minute of timestamp.
* `month` - Number month of timestamp.
* `second` - Number second of timestamp.
* `unix` - Number of seconds since epoch time, e.g. `1581489373`.
* `year` - Number year of timestamp.

## Import

This resource can be imported using the base UTC RFC3339 value and rotation years, months, days, hours, and minutes, separated by commas (`,`), e.g. for 30 days

```console
$ terraform import time_rotation.example 2020-02-12T06:36:13Z,0,0,30,0,0
```

Otherwise, to import with the rotation RFC3339 value, the base UTC RFC3339 value and rotation UTC RFC3339 value, separated by commas (`,`), e.g.

```console
$ terraform import time_rotation.example 2020-02-12T06:36:13Z,2020-02-13T06:36:13Z
```

The `keepers` argument cannot be imported.
