resource "looker_user_attribute" "my_user_attribute" {
  name  = "my_name"
  label = "Display Label"
  type  = "advanced_filter_string"
}
