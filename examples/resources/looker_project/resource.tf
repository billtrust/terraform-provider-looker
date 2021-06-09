resource "looker_git_deploy_key" "project_git_deploy_key" {
  project_id = looker_project.my_project.id
}
