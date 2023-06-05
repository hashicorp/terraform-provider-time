# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "aws_ram_resource_association" "example" {
  resource_arn       = aws_subnet.example.arn
  resource_share_arn = aws_ram_resource_share.example.arn
}

# AWS resources shared via Resource Access Manager can take a few seconds to
# propagate across AWS accounts after RAM returns a successful association.
resource "time_sleep" "ram_resource_propagation" {
  create_duration = "60s"

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