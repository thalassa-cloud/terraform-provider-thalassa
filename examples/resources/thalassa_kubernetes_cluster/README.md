# Thalassa Kubernetes Cluster

## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_thalassa"></a> [thalassa](#provider\_thalassa) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| thalassa_kubernetes_cluster.example | resource |
| thalassa_kubernetes_node_pool.example | resource |
| thalassa_subnet.example | resource |
| thalassa_vpc.example | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_availability_zones"></a> [availability\_zones](#input\_availability\_zones) | Thalassa availability zones in region | `list` | <pre>[<br/>  "nl-01a",<br/>  "nl-01b",<br/>  "nl-01c"<br/>]</pre> | no |
| <a name="input_organisation_id"></a> [organisation\_id](#input\_organisation\_id) | Thalassa organisation ID | `string` | `""` | no |
| <a name="input_region"></a> [region](#input\_region) | Thalassa region | `string` | `"nl-01"` | no |
| <a name="input_thalassa_api"></a> [thalassa\_api](#input\_thalassa\_api) | Thalassa API URL | `string` | `"https://api.thalassa.cloud"` | no |
| <a name="input_thalassa_token"></a> [thalassa\_token](#input\_thalassa\_token) | Thalassa API token | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_kubernetes_cluster_id"></a> [kubernetes\_cluster\_id](#output\_kubernetes\_cluster\_id) | Output the Kubernetes cluster details |
| <a name="output_kubernetes_cluster_name"></a> [kubernetes\_cluster\_name](#output\_kubernetes\_cluster\_name) | n/a |
| <a name="output_node_pool_id"></a> [node\_pool\_id](#output\_node\_pool\_id) | # Output the node pool details |
| <a name="output_node_pool_name"></a> [node\_pool\_name](#output\_node\_pool\_name) | n/a |
