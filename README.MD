# Terraform Provider: Looker

This is a terraform provider plugin for managing [Looker](https://www.looker.com/) accounts.

## Install

### For Terraform v0.13+ users

You can use [Explicit Provider Source Locations](https://www.terraform.io/upgrade-guides/0-13.html#explicit-provider-source-locations).

```terraform
terraform {
  required_providers {
    looker = {
      source = "DevotedHealth/looker"
      version = "0.1.0"
    }
  }
}
```

## Usage

In-depth docs are available [on the Terraform registry](https://registry.terraform.io/providers/DevotedHealth/looker/latest).
