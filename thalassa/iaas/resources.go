package iaas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var (
	ResourcesMap = map[string]*schema.Resource{
		"thalassa_block_volume":            resourceBlockVolume(),
		"thalassa_block_volume_attachment": resourceBlockVolumeAttachment(),

		"thalassa_loadbalancer_listener":    resourceLoadBalancerListener(),
		"thalassa_loadbalancer":             resourceLoadBalancer(),
		"thalassa_natgateway":               resourceNatGateway(),
		"thalassa_route_table_route":        resourceRouteTableRoute(),
		"thalassa_route_table":              resourceRouteTable(),
		"thalassa_security_group":           ResourceSecurityGroup(),
		"thalassa_subnet":                   resourceSubnet(),
		"thalassa_target_group_attachment":  resourceTargetGroupAttachment(),
		"thalassa_target_group":             resourceTargetGroup(),
		"thalassa_virtual_machine_instance": resourceVirtualMachineInstance(),
		"thalassa_vpc":                      resourceVpc(),
		"thalassa_vpc_firewall_rule":        resourceVpcFirewallRule(),
		"thalassa_cloud_init_template":      resourceCloudInitTemplate(),
	}

	DataSourcesMap = map[string]*schema.Resource{
		"thalassa_region":             DataSourceRegion(),
		"thalassa_regions":            DataSourceRegions(),
		"thalassa_machine_image":      DataSourceMachineImage(),
		"thalassa_machine_type":       DataSourceMachineType(),
		"thalassa_vpc":                DataSourceVpc(),
		"thalassa_vpc_firewall_rule":  DataSourceVpcFirewallRule(),
		"thalassa_vpc_firewall_rules": DataSourceVpcFirewallRules(),
		"thalassa_security_group":     DataSourceSecurityGroup(),
		"thalassa_volume_type":        DataSourceVolumeType(),
		"thalassa_subnet":             dataSourceSubnet(),
		"thalassa_natgateway":         DataSourceNatGateway(),
		"thalassa_loadbalancer":       DataSourceLoadBalancer(),
	}
)
