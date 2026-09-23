resource "time_rotating_v2" "example" {
  rotation_days = 30
}

output "rotation_timestamp" {
  value = time_rotating_v2.example.next_rotation_rfc3339
}
