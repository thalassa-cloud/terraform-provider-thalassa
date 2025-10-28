package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"

	iaas "github.com/thalassa-cloud/client-go/iaas"
)

func resourceVpcFirewallRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a VPC firewall rule",
		CreateContext: resourceVpcFirewallRuleCreate,
		ReadContext:   resourceVpcFirewallRuleRead,
		UpdateContext: resourceVpcFirewallRuleUpdate,
		DeleteContext: resourceVpcFirewallRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity of the VPC that the firewall rule belongs to",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 16),
				Description:  "Name of the VPC firewall rule. Must be between 1 and 16 characters and contain only ASCII characters.",
			},
			"protocols": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Protocols that the firewall rule applies to",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the rule applies to TCP protocol",
						},
						"udp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the rule applies to UDP protocol",
						},
						"icmp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the rule applies to ICMP protocol",
						},
						"any": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the rule applies to any protocol",
						},
					},
				},
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "Source CIDR of the firewall rule",
			},
			"source_ports": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Source ports of the firewall rule. Must be between 0 and 65535",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validate.IntBetween(0, 65535),
				},
			},
			"destination": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.IsCIDR,
				Description:  "Destination CIDR of the firewall rule",
			},
			"destination_ports": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Destination ports of the firewall rule. Must be between 0 and 65535",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validate.IntBetween(0, 65535),
				},
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringInSlice([]string{"allow", "drop"}, false),
				Description:  "Action of the firewall rule. One of allow, drop",
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.IntBetween(1, 1000),
				Description:  "Priority of the firewall rule. Must be between 1 and 1000",
			},
			"source_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the source subnet",
			},
			"destination_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the destination subnet",
			},
			"interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity of the interface. Leaving empty will apply to all interfaces",
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringInSlice([]string{"inbound", "outbound"}, false),
				Description:  "Direction of the firewall rule",
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "active",
				ValidateFunc: validate.StringInSlice([]string{"active", "inactive", "deleted"}, false),
				Description:  "State of the firewall rule",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the firewall rule was created",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVpcFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	protocols := d.Get("protocols").([]interface{})[0].(map[string]interface{})
	protocolsConfig := iaas.VpcFirewallRuleProtocols{
		TCP:  protocols["tcp"].(bool),
		UDP:  protocols["udp"].(bool),
		ICMP: protocols["icmp"].(bool),
		Any:  protocols["any"].(bool),
	}

	createRequest := iaas.CreateVpcFirewallRuleRequest{
		Name:        d.Get("name").(string),
		VpcIdentity: d.Get("vpc_id").(string),
		Protocols:   protocolsConfig,
		Action:      iaas.FirewallRuleAction(d.Get("action").(string)),
		Priority:    int32Ptr(int32(d.Get("priority").(int))),
		Direction:   iaas.VpcFirewallRuleDirection(d.Get("direction").(string)),
		State:       iaas.FirewallRuleState(d.Get("state").(string)),
	}

	if source, ok := d.GetOk("source"); ok {
		sourceStr := source.(string)
		createRequest.Source = &sourceStr
	}

	if destination, ok := d.GetOk("destination"); ok {
		destinationStr := destination.(string)
		createRequest.Destination = &destinationStr
	}

	if sourceSubnetIdentity, ok := d.GetOk("source_subnet_id"); ok {
		sourceSubnetStr := sourceSubnetIdentity.(string)
		createRequest.SourceSubnetIdentity = &sourceSubnetStr
	}

	if destinationSubnetIdentity, ok := d.GetOk("destination_subnet_id"); ok {
		destinationSubnetStr := destinationSubnetIdentity.(string)
		createRequest.DestinationSubnetIdentity = &destinationSubnetStr
	}

	if interfaceIdentity, ok := d.GetOk("interface_id"); ok {
		interfaceStr := interfaceIdentity.(string)
		createRequest.InterfaceIdentity = &interfaceStr
	}

	if sourcePorts, ok := d.GetOk("source_ports"); ok {
		ports := convert.ConvertToInt32Slice(sourcePorts)
		createRequest.SourcePorts = ports
	}

	if destinationPorts, ok := d.GetOk("destination_ports"); ok {
		ports := convert.ConvertToInt32Slice(destinationPorts)
		createRequest.DestinationPorts = ports
	}

	firewallRule, err := client.IaaS().CreateVpcFirewallRule(ctx, d.Get("vpc_id").(string), createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(firewallRule.Identity)
	return resourceVpcFirewallRuleRead(ctx, d, m)
}

func resourceVpcFirewallRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Extract VPC identity from the resource ID or from state
	vpcIdentity := d.Get("vpc_id").(string)
	if vpcIdentity == "" {
		// Try to get it from the state if not set
		vpcIdentity = d.Get("vpc_id").(string)
	}

	firewallRule, err := client.IaaS().GetVpcFirewallRule(ctx, vpcIdentity, d.Id())
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return setVpcFirewallRuleData(d, firewallRule)
}

func resourceVpcFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	protocols := d.Get("protocols").([]interface{})[0].(map[string]interface{})
	protocolsConfig := iaas.VpcFirewallRuleProtocols{
		TCP:  protocols["tcp"].(bool),
		UDP:  protocols["udp"].(bool),
		ICMP: protocols["icmp"].(bool),
		Any:  protocols["any"].(bool),
	}

	updateRequest := iaas.UpdateVpcFirewallRuleRequest{
		Identity:  d.Id(),
		Name:      d.Get("name").(string),
		Protocols: protocolsConfig,
		Action:    iaas.FirewallRuleAction(d.Get("action").(string)),
		Priority:  int32(d.Get("priority").(int)),
		Direction: iaas.VpcFirewallRuleDirection(d.Get("direction").(string)),
		State:     iaas.FirewallRuleState(d.Get("state").(string)),
	}

	if source, ok := d.GetOk("source"); ok {
		sourceStr := source.(string)
		updateRequest.Source = &sourceStr
	}

	if destination, ok := d.GetOk("destination"); ok {
		destinationStr := destination.(string)
		updateRequest.Destination = &destinationStr
	}

	if sourceSubnetIdentity, ok := d.GetOk("source_subnet_id"); ok {
		sourceSubnetStr := sourceSubnetIdentity.(string)
		updateRequest.SourceSubnetIdentity = &sourceSubnetStr
	}

	if destinationSubnetIdentity, ok := d.GetOk("destination_subnet_id"); ok {
		destinationSubnetStr := destinationSubnetIdentity.(string)
		updateRequest.DestinationSubnetIdentity = &destinationSubnetStr
	}

	if interfaceIdentity, ok := d.GetOk("interface_id"); ok {
		interfaceStr := interfaceIdentity.(string)
		updateRequest.InterfaceIdentity = &interfaceStr
	}

	if sourcePorts, ok := d.GetOk("source_ports"); ok {
		ports := convert.ConvertToInt32Slice(sourcePorts)
		updateRequest.SourcePorts = ports
	}

	if destinationPorts, ok := d.GetOk("destination_ports"); ok {
		ports := convert.ConvertToInt32Slice(destinationPorts)
		updateRequest.DestinationPorts = ports
	}

	firewallRule, err := client.IaaS().UpdateVpcFirewallRule(ctx, d.Get("vpc_id").(string), d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	return setVpcFirewallRuleData(d, firewallRule)
}

func resourceVpcFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.IaaS().DeleteVpcFirewallRule(ctx, d.Get("vpc_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func setVpcFirewallRuleData(d *schema.ResourceData, rule *iaas.VpcFirewallRule) diag.Diagnostics {
	d.SetId(rule.Identity)
	d.Set("name", rule.Name)
	d.Set("action", rule.Action)
	d.Set("priority", rule.Priority)
	d.Set("direction", rule.Direction)
	d.Set("state", rule.State)
	d.Set("created_at", rule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))

	if rule.Source != nil {
		d.Set("source", *rule.Source)
	}
	if rule.Destination != nil {
		d.Set("destination", *rule.Destination)
	}
	if rule.SourceSubnet != nil {
		d.Set("source_subnet_id", rule.SourceSubnet.Identity)
	}
	if rule.DestinationSubnet != nil {
		d.Set("destination_subnet_id", rule.DestinationSubnet.Identity)
	}
	if rule.Interface != nil {
		d.Set("interface_id", rule.Interface.Identity)
	}

	// Set protocols
	protocols := []map[string]interface{}{
		{
			"tcp":  rule.Protocols.TCP,
			"udp":  rule.Protocols.UDP,
			"icmp": rule.Protocols.ICMP,
			"any":  rule.Protocols.Any,
		},
	}
	d.Set("protocols", protocols)

	// Set source ports
	if rule.SourcePorts != nil {
		sourcePorts := make([]int, len(rule.SourcePorts))
		for i, port := range rule.SourcePorts {
			sourcePorts[i] = int(port)
		}
		d.Set("source_ports", sourcePorts)
	}

	// Set destination ports
	if rule.DestinationPorts != nil {
		destinationPorts := make([]int, len(rule.DestinationPorts))
		for i, port := range rule.DestinationPorts {
			destinationPorts[i] = int(port)
		}
		d.Set("destination_ports", destinationPorts)
	}

	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
