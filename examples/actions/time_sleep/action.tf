
resource "terraform_data" "trigger" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.time_sleep.sleep]
    }
  }
}

action "time_sleep" "sleep" {
  config {
    duration = "10s"
  }
}