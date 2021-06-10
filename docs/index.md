---
page_title: "Provider: Looker"
description: Manage Looker with Terraform.
---

# Looker Provider

This is a terraform provider plugin for managing [Looker](https://www.looker.com/) accounts.
Coverage is focused on the part of Looker related to access control.


## Example Provider Configuration

```terraform
provider "looker" {
  // required
  client_id     = "..."
  client_secret = "..."
  base_url      = "..."
}
```
