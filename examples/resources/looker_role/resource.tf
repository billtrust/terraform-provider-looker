resource "looker_role" "role" {
  name              = "Role"
  permission_set_id = looker_permission_set.embed_permission_set.id
  model_set_id      = looker_model_set.model_set.id
}
