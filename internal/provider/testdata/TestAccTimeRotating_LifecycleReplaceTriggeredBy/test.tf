# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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