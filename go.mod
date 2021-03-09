module github.com/billtrust/terraform-provider-looker

go 1.15

require (
	github.com/billtrust/looker-go-sdk v0.0.0-20210308030254-15099837719d
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.0
	github.com/hashicorp/terraform v0.14.6
	golang.org/dl v0.0.0-20210220033039-562909534da3 // indirect
)

replace github.com/billtrust/looker-go-sdk => ../looker-go-sdk
