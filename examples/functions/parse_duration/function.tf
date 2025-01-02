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
  value = provider::time::parse_duration("1h")
}
