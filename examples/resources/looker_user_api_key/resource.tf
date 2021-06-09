resource "looker_user_api_key" "user_api_key" {
  user_id = looker_user.user.id

  provisioner "local-exec" {
    command    = "google-chrome ${var.looker_base_url}/admin/users/api3_key/${self.user_id}"
    on_failure = "continue"
  }
}
