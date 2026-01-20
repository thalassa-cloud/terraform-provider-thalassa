package dbaas

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/convert"
	"github.com/thalassa-cloud/terraform-provider-thalassa/thalassa/provider"
)

func resourcePgRoles() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a PostgreSQL role",
		CreateContext: resourcePgRolesCreate,
		ReadContext:   resourcePgRolesRead,
		UpdateContext: resourcePgRolesUpdate,
		DeleteContext: resourcePgRolesDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the PostgreSQL role",
			},
			"organisation_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Reference to the Organisation of the Db Cluster. If not provided, the organisation of the (Terraform) provider will be used.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "The name of the role",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					//Field name may only contain lowercase alphanumeric characters & underscores
					if !regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(val.(string)) {
						errs = append(errs, fmt.Errorf("name may only contain lowercase alphanumeric characters & underscores"))
					}
					warns = []string{}
					return
				},
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The password of the role",
				Sensitive:   true,
			},
			"db_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database",
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				Description: "The connection limit of the role",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val.(int) < -1 {
						errs = append(errs, fmt.Errorf("connection_limit must be greater than -1"))
					}
					warns = []string{}
					return
				},
			},
			"create_db": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the role can create databases",
			},
			"create_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the role can create roles",
			},
			"login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the role can login",
			},
		},
	}
}

func resourcePgRolesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			return diag.FromErr(err)
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}

		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	createRole := dbaas.CreatePgRoleRequest{
		Name:            d.Get("name").(string),
		Password:        d.Get("password").(string),
		ConnectionLimit: int64(d.Get("connection_limit").(int)),
		CreateDb:        d.Get("create_db").(bool),
		CreateRole:      d.Get("create_role").(bool),
		Login:           d.Get("login").(bool),
	}

	dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if the role already exists
	for _, role := range dbCluster.PostgresRoles {
		if strings.EqualFold(role.Name, createRole.Name) {
			d.SetId(role.Identity)
			d.Set("name", role.Name)
			d.Set("db_cluster_id", dbClusterId)
			d.Set("connection_limit", role.ConnectionLimit)
			d.Set("create_db", role.CreateDb)
			d.Set("create_role", role.CreateRole)
			d.Set("login", role.Login)
			return resourcePgRolesRead(ctx, d, m)
		}
	}

	err = client.DBaaS().CreatePgRole(ctx, dbCluster.Identity, createRole)
	if err != nil {
		return diag.FromErr(err)
	}

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			return diag.FromErr(err)
		}
		found := false
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			for _, role := range dbCluster.PostgresRoles {
				if strings.EqualFold(role.Name, createRole.Name) {
					d.SetId(role.Identity)
					d.Set("name", role.Name)
					d.Set("db_cluster_id", dbClusterId)
					d.Set("connection_limit", role.ConnectionLimit)
					d.Set("create_db", role.CreateDb)
					d.Set("create_role", role.CreateRole)
					d.Set("login", role.Login)
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	return resourcePgRolesRead(ctx, d, m)
}

func resourcePgRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Get("id").(string)

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			return diag.FromErr(err)
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	for _, role := range dbCluster.PostgresRoles {
		if role.Identity == id {
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

func resourcePgRolesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil // a deleted db cluster means the pg role is also deleted
			}
			return diag.FromErr(err)
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	updateRole := dbaas.UpdatePgRoleRequest{
		ConnectionLimit: int64(d.Get("connection_limit").(int)),
	}

	if d.HasChange("password") {
		if password, ok := d.GetOk("password"); ok {
			if strVal, ok := password.(string); ok && strVal != "" {
				updateRole.Password = convert.Ptr(strVal)
			}
		}
	}

	err = client.DBaaS().UpdatePgRole(ctx, dbCluster.Identity, d.Get("id").(string), updateRole)
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg role is also deleted
		}
		return diag.FromErr(err)
	}

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			return diag.FromErr(err)
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	return resourcePgRolesRead(ctx, d, m)
}

func resourcePgRolesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := provider.GetClient(provider.GetProvider(m), d)
	if err != nil {
		return diag.FromErr(err)
	}

	dbClusterId := d.Get("db_cluster_id").(string)
	var dbCluster *dbaas.DbCluster

	for {
		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		default:
			time.Sleep(1 * time.Second)
		}
		dbCluster, err = client.DBaaS().GetDbCluster(ctx, dbClusterId)
		if err != nil {
			if tcclient.IsNotFound(err) {
				d.SetId("")
				return nil // a deleted db cluster means the pg role is also deleted
			}
			return diag.FromErr(err)
		}

		if dbCluster == nil {
			return diag.FromErr(fmt.Errorf("db cluster not found"))
		}
		if dbCluster.Status == dbaas.DbClusterStatusReady {
			break
		}
	}

	err = client.DBaaS().DeletePgRole(ctx, dbCluster.Identity, d.Get("id").(string))
	if err != nil {
		if tcclient.IsNotFound(err) {
			d.SetId("")
			return nil // a deleted db cluster means the pg role is also deleted
		}
		return diag.FromErr(err)
	}

	// TODO: Wait for the role to be deleted this is currently not correctly implemented in the API

	d.SetId("")
	return nil
}
