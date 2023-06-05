# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "time_static" "example" {}

output "current_time" {
  value = time_static.example.rfc3339
}