resource "time_offset" "example" {
  offset_days = 7
}

output "one_week_from_now" {
  value = time_offset.example.rfc3339
}