package thalassa

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_TOKEN", nil),
			},
			"api": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_ENDPOINT", "https://api.thalassa.cloud"),
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_ORGANISATION", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"thalassa_vpc":                      resourceVpc(),
			"thalassa_subnet":                   resourceSubnet(),
			"thalassa_route_table":              resourceRouteTable(),
			"thalassa_route_table_route":        resourceRouteTableRoute(),
			"thalassa_loadbalancer":             resourceLoadBalancer(),
			"thalassa_natgateway":               resourceNatGateway(),
			"thalassa_kubernetes_cluster":       resourceKubernetesCluster(),
			"thalassa_kubernetes_node_pool":     resourceKubernetesNodePool(),
			"thalassa_block_volume":             resourceBlockVolume(),
			"thalassa_virtual_machine_instance": resourceVirtualMachineInstance(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thalassa_organisation":       dataSourceOrganisations(),
			"thalassa_region":             dataSourceRegion(),
			"thalassa_regions":            dataSourceRegions(),
			"thalassa_kubernetes_version": dataSourceKubernetesVersion(),
			"thalassa_kubernetes_cluster": dataSourceKubernetesCluster(),
			"thalassa_machine_image":      dataSourceMachineImage(),
			"thalassa_machine_type":       dataSourceMachineType(),
			"thalassa_vpc":                dataSourceVpc(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func getProvider(m interface{}) ConfiguredProvider {
	p, ok := m.(ConfiguredProvider)
	if !ok {
		panic("invalid configured provider")
	}
	return p
}

type ConfiguredProvider struct {
	Client       thalassa.Client
	Organisation string
	token        string
	apiEndpoint  string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	apiEndpoint := d.Get("api").(string)
	organisation := d.Get("organisation_id").(string)

	internalClient, err := thalassa.NewClient(
		client.WithBaseURL(apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithAuthPersonalToken(token),
		client.WithUserAgent("thalassa-cloud/terraform-provider-thalassa"),
	)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return ConfiguredProvider{
		Client:       internalClient,
		Organisation: organisation,
		token:        token,
		apiEndpoint:  apiEndpoint,
	}, nil
}

func getClient(provider ConfiguredProvider, d *schema.ResourceData) (thalassa.Client, error) {
	organisation, err := getOrganisation(provider, d)
	if err != nil {
		return nil, err
	}
	client, err := thalassa.NewClient(
		client.WithBaseURL(provider.apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithAuthPersonalToken(provider.token),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getOrganisation(provider ConfiguredProvider, d *schema.ResourceData) (string, error) {
	organisation := d.Get("organisation_id").(string)
	if organisation != "" {
		return organisation, nil
	}
	if provider.Organisation != "" {
		return provider.Organisation, nil
	}
	return "", errors.New("organisation is not set")
}
