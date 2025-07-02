package dbaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func dataSourcePgRoles() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourcePgRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role",
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database cluster",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connection limit of the role",
			},
			"create_db": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role can create databases",
			},
			"create_role": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role can create roles",
			},
			"login": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role can login",
			},
		},
	}
}

func dataSourcePgRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	dbCluster, err := client.DbaaSAlphaV1().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, role := range dbCluster.PostgresRoles {
		if role.Name == d.Get("name").(string) {
			d.SetId(role.Identity)
			d.Set("name", role.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("connection_limit", role.ConnectionLimit)
			d.Set("create_db", role.CreateDb)
			d.Set("create_role", role.CreateRole)
			d.Set("login", role.Login)
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("role not found"))
}
