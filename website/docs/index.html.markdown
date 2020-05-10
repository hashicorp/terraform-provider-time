---
layout: "time"
page_title: "Provider: Time"
description: |-
  The time provider is used to interact with time-based resources.
---

# Time Provider

The time provider is used to interact with time-based resources. The provider itself has no configuration options.

Use the navigation to the left to read about the available resources.

## Resource "Triggers"

Certain time resources, only perform actions during specific lifecycle actions:

- `time_offset`: Saves base timestamp into Terraform state only when created.
- `time_sleep`: Sleeps when created and/or destroyed.
- `time_static`: Saves base timestamp into Terraform state only when created.

These resources provide an optional map argument called `triggers` that can be populated with arbitrary key/value pairs. When the keys or values of this argument are updated, Terraform will re-perform the desired action, such as updating the base timestamp or sleeping again.

For example:

```hcl
resource "time_static" "ami_update" {
  triggers = {
    # Save the time each switch of an AMI id
    ami_id = data.aws_ami.example.id
  }
}

resource "aws_instance" "server" {
  # Read the AMI id "through" the time_static resource to ensure that
  # both will change together.
  ami = time_static.ami_update.triggers.ami_id

  tags = {
    AmiUpdateTime = time_static.ami_update.rfc3339
  }

  # ... (other aws_instance arguments) ...
}
```

`triggers` are *not* treated as sensitive attributes; a value used for `triggers` will be displayed in Terraform UI output as plaintext.

To force a these actions to reoccur without updating `triggers`, the [`terraform taint` command](https://www.terraform.io/docs/commands/taint.html) can be used to produce the action on the next run.
