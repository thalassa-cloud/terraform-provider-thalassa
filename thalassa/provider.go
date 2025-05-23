package thalassa

import (
	"context"
	"errors"
	"fmt"

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
				Optional:    true,
				Sensitive:   true,
				Description: "The API token for authentication. Can be set via the THALASSA_API_TOKEN environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_TOKEN", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The OIDC client ID for authentication. Can be set via the THALASSA_CLIENT_ID environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The OIDC client secret for authentication. Can be set via the THALASSA_CLIENT_SECRET environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_CLIENT_SECRET", nil),
			},
			"api": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API endpoint URL. Can be set via the THALASSA_API_ENDPOINT environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_API_ENDPOINT", "https://api.thalassa.cloud"),
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organisation ID to use. Can be set via the THALASSA_ORGANISATION environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("THALASSA_ORGANISATION", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"thalassa_block_volume":             resourceBlockVolume(),
			"thalassa_block_volume_attachment":  resourceBlockVolumeAttachment(),
			"thalassa_kubernetes_cluster":       resourceKubernetesCluster(),
			"thalassa_kubernetes_node_pool":     resourceKubernetesNodePool(),
			"thalassa_loadbalancer_listener":    resourceLoadBalancerListener(),
			"thalassa_loadbalancer":             resourceLoadBalancer(),
			"thalassa_natgateway":               resourceNatGateway(),
			"thalassa_route_table_route":        resourceRouteTableRoute(),
			"thalassa_route_table":              resourceRouteTable(),
			"thalassa_security_group":           resourceSecurityGroup(),
			"thalassa_subnet":                   resourceSubnet(),
			"thalassa_target_group_attachment":  resourceTargetGroupAttachment(),
			"thalassa_target_group":             resourceTargetGroup(),
			"thalassa_virtual_machine_instance": resourceVirtualMachineInstance(),
			"thalassa_vpc":                      resourceVpc(),
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
			"thalassa_security_group":     dataSourceSecurityGroup(),
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
	clientID     string
	clientSecret string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	apiEndpoint := d.Get("api").(string)
	organisation := d.Get("organisation_id").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	opts := []client.Option{
		client.WithBaseURL(apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithUserAgent("thalassa-cloud/terraform-provider-thalassa"),
	}

	hasAuth := false
	if token != "" {
		opts = append(opts, client.WithAuthPersonalToken(token))
		hasAuth = true
	}

	if clientID != "" && clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(clientID, clientSecret, fmt.Sprintf("%s/oidc/token", apiEndpoint)))
		hasAuth = true
	}

	if !hasAuth {
		return nil, diag.FromErr(errors.New("no authentication method provided"))
	}

	internalClient, err := thalassa.NewClient(opts...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return ConfiguredProvider{
		Client:       internalClient,
		Organisation: organisation,
		token:        token,
		apiEndpoint:  apiEndpoint,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

func getClient(provider ConfiguredProvider, d *schema.ResourceData) (thalassa.Client, error) {
	organisation, err := getOrganisation(provider, d)
	if err != nil {
		return nil, err
	}

	opts := []client.Option{
		client.WithBaseURL(provider.apiEndpoint),
		client.WithOrganisation(organisation),
		client.WithUserAgent("thalassa-cloud/terraform-provider-thalassa"),
	}

	hasAuth := false
	if provider.token != "" {
		opts = append(opts, client.WithAuthPersonalToken(provider.token))
		hasAuth = true
	}

	if provider.clientID != "" && provider.clientSecret != "" {
		opts = append(opts, client.WithAuthOIDC(provider.clientID, provider.clientSecret, fmt.Sprintf("%s/oidc/token", provider.apiEndpoint)))
		hasAuth = true
	}

	if !hasAuth {
		return nil, errors.New("no authentication method provided")
	}

	client, err := thalassa.NewClient(opts...)
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
