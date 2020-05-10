---
layout: "time"
page_title: "Time: time_sleep"
description: |-
  Manages a static time resource.
---

# Resource: time_sleep

Manages a resource that delays creation and/or destruction, typically for further resources. This prevents cross-platform compatibility and destroy-time issues with using the [`local-exec` provisioner](https://www.terraform.io/docs/provisioners/local-exec.html).

-> In many cases, this resource should be considered a workaround for issues that should be reported and handled in downstream Terraform Provider logic. Downstream resources can usually introduce or adjust retries in their code to handle time delay issues for all Terraform configurations.

## Example Usage

### Creation Delay Usage

```hcl
# This resource will destroy (potentially immediately) after null_resource.next
resource "null_resource" "previous" {}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [null_resource.previous]

  create_seconds = 30
}

# This resource will create (at least) 30 seconds after null_resource.previous
resource "null_resource" "next" {
  depends_on = [time_sleep.wait_30_seconds]
}
```

### Destruction Delay Usage

```hcl
# This resource will destroy (at least) 30 seconds after null_resource.next
resource "null_resource" "previous" {}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [null_resource.previous]

  destroy_seconds = 30
}

# This resource will create (potentially immediately) after null_resource.previous
resource "null_resource" "next" {
  depends_on = [time_sleep.wait_30_seconds]
}
```

### Triggers Usage

```hcl
resource "aws_ram_resource_association" "example" {
  resource_arn       = aws_subnet.example.arn
  resource_share_arn = aws_ram_resource_share.example.arn
}

# AWS resources shared via Resource Access Manager can take a few seconds to
# propagate across AWS accounts after RAM returns a successful association.
resource "time_sleep" "ram_resource_propagation" {
  create_seconds = 60

  triggers = {
    # This sets up a proper dependency on the RAM association
    subnet_arn = aws_ram_resource_association.example.resource_arn
    subnet_id  = aws_subnet.example.id
  }
}

resource "aws_db_subnet_group" "example" {
  name = "example"

  # Read the Subnet identifier "through" the time_sleep resource to ensure a
  # proper dependency and that both will change together.
  subnet_ids = [time_sleep.ram_resource_propagation.triggers["subnet_id"]]
}
```

## Argument Reference

The following arguments are optional:

* `create_seconds` - (Optional) Number of seconds to sleep on resource creation. Updating this value by itself will not trigger sleeping.
* `destroy_seconds` - (Optional) Number of seconds to sleep on resource destroy. Updating this value by itself will not trigger sleeping. This value or any updates to it must be successfully applied into the Terraform state before destroying this resource to take effect.
* `triggers` - (Optional) Arbitrary map of values that, when changed, will run any creation or destroy delays again. See [the main provider documentation](../index.html) for more information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - UTC RFC3339 timestamp of the creation or import, e.g. `2020-02-12T06:36:13Z`.

## Import

This resource can be imported with the `create_seconds` and `destroy_seconds`, separated by a comma (`,`), e.g.

```console
$ terraform import time_sleep.example 30,0
```

The `triggers` argument cannot be imported.
