---
layout: "time"
page_title: "Time: time_offset"
description: |-
  Manages a offset time resource.
---

# Resource: time_offset

Manages a offset time resource, which keeps an UTC timestamp saved in the Terraform state that is offset from a base timestamp. This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).

-> Further manipulation of incoming or outgoing values can be accomplished with the [`formatdate()` function](https://www.terraform.io/docs/configuration/functions/formatdate.html) and the [`timeadd()` function](https://www.terraform.io/docs/configuration/functions/timeadd.html).

## Example Usage

```hcl
resource "time_offset" "example" {
  offset_days = 7
}

output "one_week_from_now" {
  value = time_offset.example.rfc3339
}
```

## Argument Reference

~> **NOTE:** At least one of the `offset_` arguments must be configured.

The following arguments are optional:

* `base_rfc3339` - (Optional) Configure the base timestamp with an UTC [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.
* `offset_days` - (Optional) Number of days to offset the base timestamp. Conflicts with other `offset_` arguments.
* `offset_hours` - (Optional) Number of hours to offset the base timestamp. Conflicts with other `offset_` arguments.
* `offset_minutes` - (Optional) Number of minutes to offset the base timestamp. Conflicts with other `offset_` arguments.
* `offset_months` - (Optional) Number of months to offset the base timestamp. Conflicts with other `offset_` arguments.
* `offset_seconds` - (Optional) Number of seconds to offset the base timestamp. Conflicts with other `offset_` arguments.
* `offset_years` - (Optional) Number of years to offset the base timestamp. Conflicts with other `offset_` arguments.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `day` - Number day of offset timestamp.
* `hour` - Number hour of offset timestamp.
* `id` - UTC RFC3339 base timestamp format, e.g. `2020-02-12T06:36:13Z`.
* `minute` - Number minute of offset timestamp.
* `month` - Number month of offset timestamp.
* `rfc822` - RFC822 timestamp (named timezone) format, e.g. `12 Feb 20 06:36 UTC`.
* `rfc822z` - RFC822 timestamp (+/-#### time offset) format, e.g. `12 Feb 20 06:36 +0000`.
* `rfc850` - RFC850 timestamp format, e.g. `Wednesday, 12-Feb-20 06:36:13 UTC`
* `rfc1123` - RFC1123 timestamp (named timezone) format, e.g. `Wed, 12 Feb 2020 06:36:13 UTC`.
* `rfc1123z` - RFC1123 timestamp (+/-#### time offset) format, e.g. `Wed, 12 Feb 2020 06:36:13 +0000`.
* `rfc3339` - RFC3339 timestamp format, e.g. `2020-02-12T06:36:13Z`.
* `second` - Number second of offset timestamp.
* `unix` - Number of seconds since epoch time, e.g. `1581489373`.
* `unixdate` - UNIX date format, e.g. `Wed Feb 12 06:36:13 UTC 2020`.
* `year` - Number year of offset timestamp.

## Import

This resource can be imported using the base UTC RFC3339 timestamp and offset years, months, days, hours, minutes, and seconds, separated by commas (`,`), e.g.

```console
$ terraform import time_offset.example 2020-02-12T06:36:13Z,0,0,7,0,0,0
```
