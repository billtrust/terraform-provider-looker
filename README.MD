# Looker Terraform Provider

## Resources

* **looker_user**

* **looker_user_roles**

* **looker_user_email**

* **looker_user_api_key**

* **looker_role**

* **looker_role_groups**

* **looker_permission_set**

* **looker_model_set**

* **looker_group**

* **looker_content_metadata_access** - gives access for a group to a space with a specific permission type (view, edit)
** NOTE - I think spaces still have some edge cases when modifying resources because of SpaceID in the swagger being defined as a string, but the service returning an int64

* **looker_main_space** - space configuration for a space whose parent we have not created (Example: "Embed Groups", "Users", "Shared", and "Embed Users")

* **looker_child_space** - space configuration for a space whose parent we have created

* **looker_connection** - This is mostly implemented to support the snowflake database. Google BigQuery is also supported. More work can be done to suport other database backends.

* **looker_project** - sets up a base project with just the name

* **looker_git_deploy_key** - TODO: rename to project_git_deploy_key - This creates a private/public key within Looker.  The public key can than be added to the git repository to allow looker to deploy models to it (todo: not sure if this is the correct wording)

* **looker_project_git_details** - updates an existing project with details on the git repository.  This is currently a separate resource because Looker requires the git deploy key be added to the git repository before any of the git details can be updated on the project.

  The order of events should be:
  1. Create a Looker project with minimal properties set (currently only name)
  2. Create a ssh private/public key pair
  3. Add the public ssh key to the git repository
  4. Update the Looker project with the details of the git repository

## Development

## Build

### make (linux/mac)

```shell
# build
make
# move provider to the local path expected by terraform >= 0.13
make install
```

### Linux
```shell
docker run -it --rm  -w /go/src/github.com/billtrust/terraform-provider-looker -v `pwd`:/go/src/github.com/billtrust/terraform-provider-looker golang:1.12 bash -c "go get && rm -rf bin && mkdir bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-looker"
```

### Windows
```shell
docker run -it --rm  -w /go/src/github.com/billtrust/terraform-provider-looker -v `pwd`:/go/src/github.com/billtrust/terraform-provider-looker golang:1.12 bash -c "go get && rm -rf bin && mkdir bin && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-looker.exe"
```

### Mac
```shell
docker run -it --rm  -w /go/src/github.com/billtrust/terraform-provider-looker -v `pwd`:/go/src/github.com/billtrust/terraform-provider-looker golang:1.12 bash -c "go get && rm -rf bin && mkdir bin && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-looker"
```

## Usage

Setup
```
variable "looker_client_id" {}
variable "looker_client_secret" {}
variable "looker_base_url" {}
variable "looker_api_port" {}

provider "looker" {
  client_id     = "${var.looker_client_id}"
  client_secret = "${var.looker_client_secret}"
  base_url      = "${var.looker_base_url}:${var.looker_api_port}"
}
```

```
resource "looker_user" "user" {
  first_name = "Reporting"
  last_name  = "API User"
}
```

```
resource "looker_user_attribute" "my_user_attribute" {
  name   = "my_name"
  label  = "Display Label"
  type   = "advanced_filter_string"
}
```

```
resource "looker_user_roles" "user_roles" {
  user_id    = "${looker_user.user.id}"
  role_names = ["Admin"]
}
```

```
resource "looker_user_api_key" "user_api_key" {
  user_id = "${looker_user.user.id}"

  provisioner "local-exec" {
    command    = "google-chrome ${var.looker_base_url}/admin/users/api3_key/${self.user_id}"
    on_failure = "continue"
  }
}
```

```
resource "looker_permission_set" "embed_permission_set" {
  name        = "Embed Permission Set"
  permissions = ["access_data", "download_with_limit", "schedule_look_emails", "schedule_external_look_emails", "see_user_dashboards"]
}
```

```
resource "looker_model_set" "model_set" {
  name   = "MyModelSet"
  models = ["accounts", "documents"]
}
```

```
resource "looker_role" "embed_role" {
  name              = "Embed User Role"
  permission_set_id = "${looker_permission_set.embed_permission_set.id}"
  model_set_id      = "${looker_model_set.model_set.id}"
}
```

```
resource "looker_group" "embed_group" {
  name = "My Embed Group"
}
```

```
resource "looker_role_groups" "embed_role_groups" {
  role_id   = "${looker_role.embed_role.id}"
  group_ids = ["${looker_group.embed_group.id}"]
}
```

```
resource "looker_main_space" "my_shared_space" {
  name                      = "My Shared Space"
  parent_space_name         = "Embed Groups"
  content_metadata_inherits = false
}
```

```
resource "looker_content_metadata_access" "embed_groups_space_access" {
  group_id            = "${looker_group.embed_group.id}"
  content_metadata_id = "${looker_main_space.my_shared_space.content_metadata_id}"
  permission_type     = "view"
}
```

```
resource "looker_connection" "snowflake_connection" {
  name                   = "snowflake"
  dialect_name           = "snowflake"
  host                   = "${var.snowflake_host}"
  port                   = "443"
  database               = "MY_DATABASE"
  username               = "${var.snowflake_username}"
  password               = "${var.snowflake_password}"
  schema                 = "PUBLIC"
  jdbc_additional_params = "account=${var.snowflake_account}&warehouse=LOAD_WH"
  ssl                    = true
  db_timezone            = "UTC"
  query_timezone         = "UTC"
}
```

```
resource "looker_project" "my_project" {
  name = "My_Project_Name"
}
```

```
resource "looker_git_deploy_key" "project_git_deploy_key" {
  project_id = "${looker_project.my_project.id}"
}
```

```
resource "github_repository_deploy_key" "looker_git_deploy_key" {
  title      = "${looker_git_deploy_key.project_git_deploy_key.id}"
  repository = "my-looker-repository"
  key        = "${looker_git_deploy_key.project_git_deploy_key.ssh_deploy_key}"
  read_only  = "false"
}

```

