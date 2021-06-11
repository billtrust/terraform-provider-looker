resource "looker_role_groups" "role_groups" {
  role_id   = looker_role.role.id
  group_ids = [looker_group.group.id]
}
