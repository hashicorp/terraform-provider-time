---
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

```terraform
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

## Implementing Rotation Windows

The following example demonstrates how to use `time_rotating` to implement a rotation mechanism for tokens, ensuring overlapping availability during transitions.

### Example Usage

```terraform
resource "time_rotating" "rotate_token_1" {
  rotation_minutes = var.expiration_time_in_minutes
}

resource "time_rotating" "rotate_token_2" {
  rfc3339 = time_rotating.rotate_token_1.rotation_rfc3339
  
  # Shorter rotation window to overlap with the first token
  rotation_minutes = var.expiration_time_in_minutes / 2

  lifecycle {
    ignore_changes = [
      rfc3339
    ]
  }
}

# Replace with a token resource from another provider
resource "token_resource" "token_1" {
  expires_at = timeadd(time_rotating.rotate_token_1.rfc3339, "${var.expiration_time_in_minutes * 1.5}m")

  # ... (other token arguments) ...
}

# Replace with a token resource from another provider
resource "token_resource" "token_2" {
  expires_at = timeadd(time_rotating.rotate_token_2.rfc3339, "${var.expiration_time_in_minutes * 1.5}m")

  # ... (other token arguments) ...
}

locals {
  use_token1 = timecmp(time_rotating.rotate_token_1.rfc3339, time_rotating.rotate_token_2.rfc3339) > 0
}

output "active_token" {
  value     = local.use_token1 ? token_resource.token_1.value : token_resource.token_2.value
  sensitive = true
}
```

### Key Considerations

1. **Overlapping Availability**:
   - Token 1 and Token 2 have overlapping rotation windows to ensure seamless availability during transitions. This minimizes potential downtime or lapses in functionality.

2. **Simplified Token Resources**:
   - The `token_resource` represents a generalized token configuration. Replace this with your specific token implementation.

3. **Active Token Logic**:
   - The `local.use_token1` logic determines which token is currently active based on the rotation timestamps.

4. **Customizable Rotation**:
   - The rotation intervals (`rotation_minutes`) can be tailored to meet your system's requirements for availability and security.