# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# The "time_rotating" resource predates the `replace_triggered_by`
# lifestyle argument introduced in Terraform v1.2.0. The `replace_triggered_by` argument looks
# for an update or replacement of the supplied resource instance. Because the "time_rotating" rotation
# checking logic is run during ReadResource() and the resource is removed from state,
# a rotation is considered to be a creation of a new resource rather than an update or replacement.
# Ref: https://github.com/hashicorp/terraform-provider-time/issues/118
resource "time_rotating" "computed_rotation" {
  rotation_minutes = 1
}

resource "time_static" "static_time" {

}

resource "time_offset" "offset_time" {
  offset_minutes = 1
}

resource "time_rotating" "configured_rfc3339" {
  rfc3339 = time_static.static_time.rfc3339
  rotation_minutes = 1
}

resource "time_rotating" "configured_rotationrfc3339" {
  rotation_rfc3339 = time_offset.offset_time.rfc3339
  rotation_minutes = 1
}


resource "terraform_data" "test_computed_rotation" {
  lifecycle {
    replace_triggered_by = [
      time_rotating.computed_rotation
    ]
  }
}

resource "terraform_data" "test_configured_rfc3339" {
  lifecycle {
    replace_triggered_by = [
      time_rotating.configured_rfc3339
    ]
  }
}

resource "terraform_data" "test_configured_rotationrfc3339" {
  lifecycle {
    replace_triggered_by = [
      time_rotating.configured_rotationrfc3339
    ]
  }
}