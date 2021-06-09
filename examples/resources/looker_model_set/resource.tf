resource "looker_model_set" "model_set" {
  name   = "MyModelSet"
  models = ["accounts", "documents"]
}
