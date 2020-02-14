---
layout: "time"
page_title: "Provider: Time"
description: |-
  The time provider is used to interact with time-based resources.
---

# Time Provider

The time provider is used to interact with time-based resources. The provider itself has no configuration options.

Use the navigation to the left to read about the available resources.

## Resource "Keepers"

The time resources, except `time_rotating`, save base timestamps only when they are created; the results produced are stored in the Terraform state and re-used until the inputs change, prompting the resource to be recreated.

The resources all provide a map argument called `keepers` that can be populated with arbitrary key/value pairs that should be selected such that they remain the same until new time values are desired.

For example:

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

Resource "keepers" are optional. The other arguments to each resource must *also* remain constant in order to retain the same time result.

`keepers` are *not* treated as sensitive attributes; a value used for `keepers` will be displayed in Terraform UI output as plaintext.

To force a base timestamp to be replaced, the [`terraform taint` command](https://www.terraform.io/docs/commands/taint.html) can be used to produce a new time on the next run.
