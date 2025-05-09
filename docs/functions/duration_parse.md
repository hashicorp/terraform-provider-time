---
page_title: "duration_parse function - terraform-provider-time"
subcategory: ""
description: |-
  Parse a Go duration string https://pkg.go.dev/time#ParseDuration into an object
---

# function: duration_parse

Given a [Go duration string](https://pkg.go.dev/time#ParseDuration), will parse and return an object representation of that duration.

## Example Usage

```terraform
# Configuration using provider functions must include required_providers configuration.
terraform {
  required_providers {
    time = {
      source = "hashicorp/time"
      # Setting the provider version is a strongly recommended practice
      # version = "..."
    }
  }
  # Provider functions require Terraform 1.8 and later.
  required_version = ">= 1.8.0"
}

output "example_output" {
  value = provider::time::duration_parse("1h")
}
```

## Signature

<!-- signature generated by tfplugindocs -->
```text
duration_parse(duration string) object
```

## Arguments

<!-- arguments generated by tfplugindocs -->
1. `duration` (String) Go time package duration string to parse


## Return Type

The `object` returned from `duration_parse` has the following attributes:
- `hours` (Number) The duration as a floating point number of hours.
- `minutes` (Number) The duration as a floating point number of minutes.
- `seconds` (Number) The duration as a floating point number of seconds.
- `milliseconds` (Number) The duration as an integer number of milliseconds.
- `microseconds` (Number) The duration as an integer number of microseconds.
- `nanoseconds` (Number) The duration as an integer number of nanoseconds.
