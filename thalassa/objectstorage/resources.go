package objectstorage

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ResourcesMap = map[string]*schema.Resource{
	"thalassa_objectstorage_bucket":           resourceBucket(),
	"thalassa_objectstorage_bucket_lifecycle": resourceBucketLifecycle(),
}

var DataSourcesMap = map[string]*schema.Resource{
	"thalassa_objectstorage_bucket": DataSourceBucket(),
}
