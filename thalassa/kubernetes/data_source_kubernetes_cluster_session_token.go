package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourceKubernetesClusterSessionToken() *schema.Resource {
	return &schema.Resource{
		Description: "Get a temporary session token for authenticating with a Kubernetes cluster for the currently authenticated user or system account. This token can be used to access the cluster via kubectl or other Kubernetes tools.",
		ReadContext: dataSourceKubernetesClusterSessionTokenRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Kubernetes cluster to get a session token for",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username for cluster authentication",
			},
			"api_server_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the Kubernetes API server",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CA certificate of the API server",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The authentication token for the session",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The complete kubeconfig file content for the session",
			},
		},
	}
}

func dataSourceKubernetesClusterSessionTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID := d.Get("cluster_id").(string)
	sessionToken, err := client.Kubernetes().GetKubernetesClusterKubeconfig(ctx, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	if sessionToken == nil {
		return diag.Errorf("no session token received for cluster %s", clusterID)
	}

	d.SetId(sessionToken.Identity)
	d.Set("username", sessionToken.Username)
	d.Set("api_server_url", sessionToken.APIServerURL)
	d.Set("ca_certificate", sessionToken.CACertificate)
	d.Set("token", sessionToken.Token)
	d.Set("kubeconfig", sessionToken.Kubeconfig)

	return nil
}
