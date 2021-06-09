resource "looker_user_roles" "user_roles" {
  user_id    = looker_user.user.id
  role_names = ["Admin"]
}
