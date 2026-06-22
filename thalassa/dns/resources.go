package dns

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ResourcesMap = map[string]*schema.Resource{
	"thalassa_dns_zone":        ResourceDnsZone(),
	"thalassa_dns_record":      ResourceDnsRecord(),
	"thalassa_dns_zone_dnssec": ResourceDnsZoneDnssec(),
}

var DataSourcesMap = map[string]*schema.Resource{}
