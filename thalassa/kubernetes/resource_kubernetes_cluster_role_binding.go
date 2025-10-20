package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	kubernetes "github.com/thalassa-cloud/client-go/kubernetes"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourceKubernetesClusterRoleBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Kubernetes cluster role binding for IAM access control. This resource allows you to bind users, teams, or service accounts to Kubernetes cluster roles.",
		CreateContext: resourceKubernetesClusterRoleBindingCreate,
		ReadContext:   resourceKubernetesClusterRoleBindingRead,
		UpdateContext: resourceKubernetesClusterRoleBindingUpdate,
		DeleteContext: resourceKubernetesClusterRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the Kubernetes cluster role binding",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Kubernetes Cluster Role Binding. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.StringLenBetween(1, 255),
				Description:  "The name of the Kubernetes cluster role binding",
			},
			"slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug of the Kubernetes cluster role binding",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.StringLenBetween(0, 1000),
				Description:  "A human-readable description of the Kubernetes cluster role binding",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels for the Kubernetes cluster role binding",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Annotations for the Kubernetes cluster role binding",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Kubernetes cluster role to bind",
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the user to bind to the cluster role",
			},
			"team_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the team to bind to the cluster role",
			},
			"service_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the service account to bind to the cluster role",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-readable note for the binding",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Kubernetes cluster role binding was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the Kubernetes cluster role binding was last updated",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKubernetesClusterRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterRoleID := d.Get("cluster_role_id").(string)
	userID := d.Get("user_id").(string)
	teamID := d.Get("team_id").(string)
	serviceAccountID := d.Get("service_account_id").(string)

	// Validate that exactly one of user_id, team_id, or service_account_id is provided
	bindingCount := 0
	if userID != "" {
		bindingCount++
	}
	if teamID != "" {
		bindingCount++
	}
	if serviceAccountID != "" {
		bindingCount++
	}

	if bindingCount != 1 {
		return diag.Errorf("exactly one of user_id, team_id, or service_account_id must be provided")
	}

	createRequest := kubernetes.CreateKubernetesClusterRoleBinding{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      convert.ConvertToMap(d.Get("labels")),
		Annotations: convert.ConvertToMap(d.Get("annotations")),
	}

	// Set the appropriate identity based on what was provided
	if userID != "" {
		createRequest.UserIdentity = &userID
	} else if teamID != "" {
		createRequest.TeamIdentity = &teamID
	} else if serviceAccountID != "" {
		createRequest.ServiceAccountIdentity = &serviceAccountID
	}

	binding, err := client.Kubernetes().CreateClusterRoleBinding(ctx, clusterRoleID, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	if binding != nil {
		d.SetId(binding.Identity)
		d.Set("slug", binding.Slug)
		d.Set("created_at", binding.CreatedAt.Format(time.RFC3339))
		if binding.UpdatedAt != nil {
			d.Set("updated_at", binding.UpdatedAt.Format(time.RFC3339))
		}
	}

	return resourceKubernetesClusterRoleBindingRead(ctx, d, m)
}

func resourceKubernetesClusterRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterRoleID := d.Get("cluster_role_id").(string)
	bindingID := d.Get("id").(string)

	// Get all bindings for the cluster role
	bindings, err := client.Kubernetes().ListClusterRoleBindings(ctx, clusterRoleID)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error getting kubernetes cluster role bindings: %s", err))
	}

	// Find the specific binding
	var binding *kubernetes.KubernetesClusterRoleBinding
	for _, b := range bindings {
		if b.Identity == bindingID {
			binding = &b
			break
		}
	}

	if binding == nil {
		d.SetId("")
		return nil
	}

	d.SetId(binding.Identity)
	d.Set("name", binding.Name)
	d.Set("slug", binding.Slug)
	d.Set("description", binding.Description)
	d.Set("labels", binding.Labels)
	d.Set("annotations", binding.Annotations)
	d.Set("cluster_role_id", clusterRoleID)
	d.Set("created_at", binding.CreatedAt.Format(time.RFC3339))
	if binding.UpdatedAt != nil {
		d.Set("updated_at", binding.UpdatedAt.Format(time.RFC3339))
	}
	if binding.Note != nil {
		d.Set("note", *binding.Note)
	}

	// Set the appropriate identity field based on what's bound
	if binding.User != nil {
		d.Set("user_id", binding.User.Subject)
	} else if binding.OrganisationTeam != nil {
		d.Set("team_id", binding.OrganisationTeam.Identity)
	} else if binding.ServiceAccount != nil {
		d.Set("service_account_id", binding.ServiceAccount.Identity)
	}

	return nil
}

func resourceKubernetesClusterRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: The API doesn't seem to have an update method for cluster role bindings
	// This would need to be implemented if the API supports it
	return diag.Errorf("updating Kubernetes cluster role bindings is not currently supported")
}

func resourceKubernetesClusterRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterRoleID := d.Get("cluster_role_id").(string)
	bindingID := d.Get("id").(string)

	err = client.Kubernetes().DeleteClusterRoleBinding(ctx, clusterRoleID, bindingID)
	if err != nil {
		if !tcclient.IsNotFound(err) {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}
