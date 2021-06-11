resource "looker_permission_set" "permission_set" {
  name        = "Permission Set"
  permissions = ["access_data", "download_with_limit", "schedule_look_emails", "schedule_external_look_emails", "see_user_dashboards"]
}
