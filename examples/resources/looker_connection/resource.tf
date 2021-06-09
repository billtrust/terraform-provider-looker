resource "looker_connection" "snowflake_connection" {
  name                   = "snowflake"
  dialect_name           = "snowflake"
  host                   = var.snowflake_host
  port                   = "443"
  database               = "MY_DATABASE"
  username               = var.snowflake_username
  password               = var.snowflake_password
  schema                 = "PUBLIC"
  jdbc_additional_params = "account=${var.snowflake_account}&warehouse=LOAD_WH"
  ssl                    = true
  db_timezone            = "UTC"
  query_timezone         = "UTC"
}
