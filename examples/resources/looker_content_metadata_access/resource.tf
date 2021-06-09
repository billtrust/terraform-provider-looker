resource "looker_content_metadata_access" "embed_groups_space_access" {
  group_id            = looker_group.embed_group.id
  content_metadata_id = looker_main_space.my_shared_space.content_metadata_id
  permission_type     = "view"
}
