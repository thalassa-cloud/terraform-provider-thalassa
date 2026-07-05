# Terraform Provider Thalassa Cloud

Official [Terraform](https://www.terraform.io/) provider for [Thalassa Cloud](https://thalassa.cloud). Manage infrastructure, Kubernetes clusters, databases, identity, DNS, object storage, KMS keys, and secrets as code.

## Documentation

- Full resource and data source reference:: [registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs](https://registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs)
- Product documentation: [docs.thalassa.cloud](https://docs.thalassa.cloud/)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.0
- A Thalassa Cloud account with API access

## Quick start

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    thalassa = {
      source  = "thalassa-cloud/thalassa"
      version = "~> 0.24"
    }
  }
}

provider "thalassa" {
  organisation_id = "my-org"
}
```

Authenticate with a personal API token:

```bash
export THALASSA_API_TOKEN="your-token"
```

Then run `terraform init` and `terraform apply` as usual.

## Provider configuration

The provider block supports the following arguments. Each can also be set via environment variable.

| Argument | Environment variable | Description |
|----------|---------------------|-------------|
| `token` | `THALASSA_API_TOKEN` | Personal API token |
| `access_token` | `THALASSA_ACCESS_TOKEN` | Access token (OIDC flow) |
| `client_id` | `THALASSA_CLIENT_ID` | OIDC client ID |
| `client_secret` | `THALASSA_CLIENT_SECRET` | OIDC client secret |
| `allow_insecure_oidc` | `THALASSA_ALLOW_INSECURE_OIDC` | Allow insecure OIDC (default: `false`) |
| `api` | `THALASSA_API_ENDPOINT` | API endpoint (default: `https://api.thalassa.cloud`) |
| `organisation_id` | `THALASSA_ORGANISATION` | Default organisation ID |
| `project_id` | `THALASSA_PROJECT_ID` | Default project ID |

Many resources accept an optional `organisation_id` attribute to override the provider default.

## Supported services

### Infrastructure (IaaS)

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_vpc` | Virtual private cloud |
| `thalassa_subnet` | Subnet within a VPC |
| `thalassa_security_group` | Security group |
| `thalassa_security_group_ingress_rule` | Ingress rule |
| `thalassa_security_group_egress_rule` | Egress rule |
| `thalassa_route_table` | Route table |
| `thalassa_route_table_route` | Route within a route table |
| `thalassa_natgateway` | NAT gateway |
| `thalassa_reserved_ip` | Reserved public IP |
| `thalassa_vpc_peering_connection` | VPC peering request |
| `thalassa_vpc_peering_connection_acceptance` | Accept a VPC peering request |
| `thalassa_vpc_firewall_rule` | VPC firewall rule |
| `thalassa_virtual_machine_instance` | Virtual machine |
| `thalassa_block_volume` | Block storage volume |
| `thalassa_block_volume_attachment` | Volume attachment |
| `thalassa_snapshot` | Volume snapshot |
| `thalassa_snapshot_policy` | Automated snapshot policy |
| `thalassa_cloud_init_template` | Cloud-init template |
| `thalassa_loadbalancer` | Load balancer |
| `thalassa_loadbalancer_listener` | Load balancer listener |
| `thalassa_target_group` | Target group |
| `thalassa_target_group_attachment` | Target group attachment |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_region` | Single region |
| `thalassa_regions` | List of regions |
| `thalassa_machine_image` | Machine image |
| `thalassa_machine_type` | Machine type |
| `thalassa_volume_type` | Block volume type |
| `thalassa_vpc` | VPC |
| `thalassa_subnet` | Subnet |
| `thalassa_vpc_default_route_table` | Default route table for a VPC |
| `thalassa_route_table` | Route table |
| `thalassa_security_group` | Security group |
| `thalassa_natgateway` | NAT gateway |
| `thalassa_loadbalancer` | Load balancer |
| `thalassa_snapshot` | Snapshot |
| `thalassa_snapshot_policy` | Snapshot policy |
| `thalassa_vpc_peering_connection` | VPC peering connection |
| `thalassa_vpc_peering_connections` | VPC peering connections |
| `thalassa_vpc_firewall_rule` | VPC firewall rule |
| `thalassa_vpc_firewall_rules` | VPC firewall rules |

### Kubernetes (KaaS)

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_kubernetes_cluster` | Kubernetes cluster |
| `thalassa_kubernetes_node_pool` | Node pool |
| `thalassa_kubernetes_cluster_role` | Cluster-scoped RBAC role |
| `thalassa_kubernetes_cluster_role_binding` | Cluster-scoped RBAC role binding |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_kubernetes_version` | Available Kubernetes versions |
| `thalassa_kubernetes_cluster` | Cluster details |
| `thalassa_kubernetes_cluster_session_token` | Cluster session token |
| `thalassa_kubernetes_cluster_role` | Cluster role |

### Database as a Service (DBaaS)

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_dbaas_db_cluster` | Database cluster |
| `thalassa_dbaas_pg_database` | PostgreSQL database |
| `thalassa_dbaas_pg_roles` | PostgreSQL roles |
| `thalassa_dbaas_pg_grant` | PostgreSQL grants |
| `thalassa_dbaas_db_backupschedule` | Backup schedule |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_dbaas_db_cluster` | Database cluster |
| `thalassa_dbaas_pg_database` | PostgreSQL database |
| `thalassa_dbaas_pg_roles` | PostgreSQL roles |
| `thalassa_dbaas_db_backupschedule` | Backup schedule |
| `thalassa_dbaas_db_backup` | Backup |

### Identity & Access Management (IAM)

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_iam_team` | Team |
| `thalassa_iam_role` | IAM role |
| `thalassa_iam_role_rule` | IAM role rule |
| `thalassa_iam_role_binding` | Role binding |
| `thalassa_iam_service_account` | Service account |
| `thalassa_iam_service_account_access_credential` | Service account credential |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_iam_team` | Team |
| `thalassa_iam_role` | IAM role |
| `thalassa_iam_service_account` | Service account |
| `thalassa_iam_organisation_members` | Organisation members |

### Key Management Service (KMS)

> KMS is in early access. The API and resource schema may change.

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_kms_key` | Encryption key |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_kms_key` | KMS key |
| `thalassa_kms_summary` | KMS availability summary |

### Secrets Manager

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_secret` | Secret |
| `thalassa_secret_version` | Secret version |
| `thalassa_secret_access_policy` | Secret access policy |

### DNS

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_dns_zone` | DNS zone |
| `thalassa_dns_record` | DNS record |
| `thalassa_dns_zone_dnssec` | DNSSEC configuration |

### Object Storage

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_objectstorage_bucket` | Object storage bucket |
| `thalassa_objectstorage_bucket_lifecycle` | Bucket lifecycle rules |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_objectstorage_bucket` | Object storage bucket |

### Thalassa File Service (TFS)

**Resources**

| Resource | Description |
|----------|-------------|
| `thalassa_tfs_instance` | TFS instance |

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_tfs_instance` | TFS instance |

### Organisation

**Data sources**

| Data source | Description |
|-------------|-------------|
| `thalassa_organisation` | Organisation |

## Examples

Runnable examples live under [`examples/`](./examples/).

### Provider configuration

- [Basic provider setup](./examples/provider/) — authentication and organisation lookup

### Infrastructure (IaaS)

#### Networking

- [VPC](./examples/resources/thalassa_vpc/)
- [Subnet](./examples/resources/thalassa_subnet/)
- [Security group](./examples/resources/thalassa_security_group/)
- [Security group ingress rule](./examples/resources/thalassa_security_group_ingress_rule/)
- [Security group egress rule](./examples/resources/thalassa_security_group_egress_rule/)
- [Route table](./examples/resources/thalassa_route_table/)
- [Route table route](./examples/resources/thalassa_route_table_route/)
- [NAT gateway](./examples/resources/thalassa_natgateway/)
- [Reserved IP](./examples/resources/thalassa_reserved_ip/)
- [VPC peering connection](./examples/resources/thalassa_vpc_peering_connection/)
- [VPC peering connection acceptance](./examples/resources/thalassa_vpc_peering_connection_acceptance/)
- [VPC firewall rule](./examples/resources/thalassa_vpc_firewall_rule/)

#### Compute & storage

- [Virtual machine instance](./examples/resources/thalassa_virtual_machine_instance/)
- [Block volume](./examples/resources/thalassa_block_volume/)
- [Block volume attachment](./examples/resources/thalassa_block_volume_attachment/)
- [Snapshot](./examples/resources/thalassa_snapshot/)
- [Snapshot policy](./examples/resources/thalassa_snapshot_policy/)
- [Cloud-init template](./examples/resources/thalassa_cloud_init_template/)

#### Load balancing

- [Load balancer](./examples/resources/thalassa_loadbalancer/)
- [Load balancer listener](./examples/resources/thalassa_loadbalancer_listener/)
- [Target group](./examples/resources/thalassa_target_group/)
- [Target group attachment](./examples/resources/thalassa_target_group_attachment/)

### Kubernetes (KaaS)

- [Kubernetes cluster](./examples/resources/thalassa_kubernetes_cluster/)
- [Kubernetes node pool](./examples/resources/thalassa_kubernetes_node_pool/)
- [Kubernetes cluster role](./examples/resources/thalassa_kubernetes_cluster_role/)
- [Kubernetes cluster role binding](./examples/resources/thalassa_kubernetes_cluster_role_binding/)

### Database as a Service (DBaaS)

- [Database cluster](./examples/resources/thalassa_dbaas_db_cluster/)
- [PostgreSQL database](./examples/resources/thalassa_dbaas_pg_database/)
- [PostgreSQL roles](./examples/resources/thalassa_dbaas_pg_roles/)
- [Database backup schedule](./examples/resources/thalassa_dbaas_db_backupschedule/)

### Identity & Access Management (IAM)

- [IAM team](./examples/resources/thalassa_iam_team/)
- [IAM role](./examples/resources/thalassa_iam_role/)
- [IAM role binding](./examples/resources/thalassa_iam_role_binding/)
- [Service account](./examples/resources/thalassa_iam_service_account/)
- [Service account access credential](./examples/resources/thalassa_iam_service_account_access_credential/)

### Key Management Service (KMS)

- [KMS key](./examples/resources/thalassa_kms_key/)
- [KMS summary data source](./examples/data-sources/thalassa_kms_summary/)
- [KMS key data source](./examples/data-sources/thalassa_kms_key/)

### Secrets Manager

- [Secret](./examples/resources/thalassa_secret/)
- [Secret version](./examples/resources/thalassa_secret_version/)
- [Secret access policy](./examples/resources/thalassa_secret_access_policy/)

### DNS

- [DNS zone](./examples/resources/thalassa_dns_zone/)
- [DNS record](./examples/resources/thalassa_dns_record/)
- [DNS zone DNSSEC](./examples/resources/thalassa_dns_zone_dnssec/)

### Object Storage

- [Object storage bucket](./examples/resources/thalassa_objectstorage_bucket/)
- [Bucket lifecycle rules](./examples/resources/thalassa_objectstorage_bucket_lifecycle/)

## Development

### Prerequisites

- [Go](https://go.dev/dl/) (see `go.mod` for the required version)
- Make

### Build and install locally

Build the provider binary and install it into the local Terraform plugin directory:

```bash
make install NAMESPACE=thalassa HOSTNAME=thalassa.cloud OS_ARCH=darwin_arm64
```

Adjust `OS_ARCH` for your platform (e.g. `linux_amd64`, `darwin_amd64`).

Reference the locally installed provider in your Terraform configuration:

```hcl
terraform {
  required_providers {
    thalassa = {
      source = "thalassa.cloud/thalassa/thalassa"
    }
  }
}
```

### Run tests

#### Unit tests

Unit tests run against mocked or in-memory state:

```bash
make test
```

#### Acceptance tests

Acceptance tests create and destroy real resources in your Thalassa organisation. They are skipped unless explicitly enabled.

**Enable acceptance tests**

Set `TF_ACC=1` (or use the `make testacc` target, which sets it for you).

**Authentication**

Acceptance test configs use an empty provider block (`provider "thalassa" {}`), so credentials and defaults must come from environment variables—the same ones used in normal Terraform runs.

Required:

| Variable | Description |
|----------|-------------|
| `THALASSA_ORGANISATION` | Organisation ID or slug for the test account |
| One of the auth options below | See examples |

Auth options (pick one):

| Variables | Description |
|-----------|-------------|
| `THALASSA_API_TOKEN` | Personal API token |
| `THALASSA_ACCESS_TOKEN` | Access token |
| `THALASSA_CLIENT_ID` + `THALASSA_CLIENT_SECRET` | OIDC client credentials |

Optional:

| Variable | Description |
|----------|-------------|
| `THALASSA_TEST_REGION` | Region for resources created during tests (default: `nl-01`) |
| `THALASSA_API_ENDPOINT` | API endpoint (default: `https://api.thalassa.cloud`) |
| `THALASSA_PROJECT_ID` | Project ID when resources require one |
| `THALASSA_ALLOW_INSECURE_OIDC` | Set to `true` when using client credentials against a non-production OIDC endpoint |

Example using a personal API token:

```bash
export TF_ACC=1
export THALASSA_API_TOKEN="your-token"
export THALASSA_ORGANISATION="my-org"
export THALASSA_TEST_REGION="nl-01"   # optional

make testacc
```

Example using an access token:

```bash
export TF_ACC=1
export THALASSA_ACCESS_TOKEN="your-access-token"
export THALASSA_ORGANISATION="my-org"

make testacc
```

Example using OIDC client credentials:

```bash
export TF_ACC=1
export THALASSA_CLIENT_ID="your-client-id"
export THALASSA_CLIENT_SECRET="your-client-secret"
export THALASSA_ORGANISATION="my-org"
export THALASSA_TEST_REGION="nl-01"   # optional

make testacc
```

**Run a single acceptance test**

Pass `-run` via `TESTARGS`:

```bash
export TF_ACC=1
export THALASSA_API_TOKEN="your-token"
export THALASSA_ORGANISATION="my-org"

make testacc TESTARGS='-run TestAccVpc_basic ./thalassa/iaas/...'
```

Or invoke `go test` directly:

```bash
TF_ACC=1 \
  THALASSA_API_TOKEN="your-token" \
  THALASSA_ORGANISATION="my-org" \
  go test ./thalassa/iaas/... -v -run TestAccVpc_basic -timeout 120m
```

Acceptance tests can take several minutes per case because they wait for cloud resources to reach a ready state. The `testacc` target uses a 120-minute timeout.

**Notes**

- Use a dedicated organisation and project for acceptance testing; tests create resources with names prefixed `tf-acc-`.

### Generate documentation

Provider docs in [`docs/`](./docs/) are generated from schema definitions:

```bash
make docs
```

## Contributing

Contributions are welcome. Please open an issue or pull request on GitHub.

When adding or changing resources, run `make docs` to regenerate the registry documentation and add an example under [`examples/`](./examples/) where practical.

## License

[Apache 2.0 License](./LICENSE)
