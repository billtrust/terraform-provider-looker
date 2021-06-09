resource "looker_main_space" "my_shared_space" {
  name                      = "My Shared Space"
  parent_space_name         = "Embed Groups"
  content_metadata_inherits = false
}
