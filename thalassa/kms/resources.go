package kms

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ResourcesMap = map[string]*schema.Resource{
	"thalassa_kms_key": ResourceKmsKey(),
}

var DataSourcesMap = map[string]*schema.Resource{
	"thalassa_kms_key": DataSourceKmsKey(),
}
