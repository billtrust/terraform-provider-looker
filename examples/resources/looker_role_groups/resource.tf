resource "looker_role_groups" "role_groups" {
  role_id   = looker_role.embed_role.id
  group_ids = ["${looker_group.embed_group.id}"]
}
