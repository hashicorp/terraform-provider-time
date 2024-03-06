terraform {
  required_providers {
    time = {
      source = "hashicorp/time"
    }
  }
}

output "example_output" {
  value = provider :: time :: rfc3339_parse("2023-07-25T23:43:16Z")
}
