---
layout: "time"
page_title: "Time: time_static"
description: |-
  Manages a static time resource.
---

# Resource: time_static

Manages a static time resource, which keeps an UTC timestamp saved in the Terraform state. This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).

-> Further manipulation of incoming or outgoing values can be accomplished with the [`formatdate()` function](https://www.terraform.io/docs/configuration/functions/formatdate.html) and the [`timeadd()` function](https://www.terraform.io/docs/configuration/functions/timeadd.html).

## Example Usage

### Basic Usage

```hcl
resource "time_static" "example" {}

output "current_time" {
  value = time_static.example.rfc3339
}
```

### Keepers Usage

```hcl
resource "time_static" "ami_update" {
  keepers = {
    # Save the time each switch of an AMI id
    ami_id = data.aws_ami.example.id
  }
}

resource "aws_instance" "server" {
  # Read the AMI id "through" the time_static resource to ensure that
  # both will change together.
  ami = time_static.ami_update.keepers.ami_id

  tags = {
    AmiUpdateTime = time_static.ami_update.rfc3339
  }

  # ... (other aws_instance arguments) ...
}
```

## Argument Reference

The following arguments are optional:

* `keepers` - (Optional) Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. See [the main provider documentation](../index.html) for more information.
* `rfc3339` - (Optional) Configure the base timestamp with an UTC [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) (`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `day` - Number day of timestamp.
* `hour` - Number hour of timestamp.
* `id` - UTC RFC3339 timestamp format, e.g. `2020-02-12T06:36:13Z`.
* `minute` - Number minute of timestamp.
* `month` - Number month of timestamp.
* `rfc822` - RFC822 timestamp (named timezone) format, e.g. `12 Feb 20 06:36 UTC`.
* `rfc822z` - RFC822 timestamp (+/-#### time offset) format, e.g. `12 Feb 20 06:36 +0000`.
* `rfc850` - RFC850 timestamp format, e.g. `Wednesday, 12-Feb-20 06:36:13 UTC`
* `rfc1123` - RFC1123 timestamp (named timezone) format, e.g. `Wed, 12 Feb 2020 06:36:13 UTC`.
* `rfc1123z` - RFC1123 timestamp (+/-#### time offset) format, e.g. `Wed, 12 Feb 2020 06:36:13 +0000`.
* `rfc3339` - RFC3339 timestamp format, e.g. `2020-02-12T06:36:13Z`.
* `second` - Number second of timestamp.
* `unix` - Number of seconds since epoch time, e.g. `1581489373`.
* `unixdate` - UNIX date format, e.g. `Wed Feb 12 06:36:13 UTC 2020`.
* `year` - Number year of timestamp.

## Import

This resource can be imported using the UTC RFC3339 value, e.g.

```console
$ terraform import time_static.example 2020-02-12T06:36:13Z
```

The `keepers` argument cannot be imported.
