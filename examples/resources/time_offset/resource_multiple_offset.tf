resource "time_offset" "example" {
  offset_years  = 1
  offset_months = 1
}

output "one_year_and_month_from_now" {
  value = time_offset.example.rfc3339
}