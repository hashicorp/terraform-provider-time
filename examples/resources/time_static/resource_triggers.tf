# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

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