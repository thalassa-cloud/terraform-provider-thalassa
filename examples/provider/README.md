# Example Provider

This example provides an example on how to use the Terraform Provider

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_thalassa"></a> [thalassa](#provider\_thalassa) | 0.1.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| thalassa_route_table.route_table | resource |
| thalassa_route_table_route.route | resource |
| thalassa_subnet.subnet | resource |
| thalassa_subnet.subnet2 | resource |
| thalassa_vpc.this | resource |
| thalassa_organisation.this | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_thalassa_api"></a> [thalassa\_api](#input\_thalassa\_api) | n/a | `any` | n/a | yes |
| <a name="input_thalassa_token"></a> [thalassa\_token](#input\_thalassa\_token) | n/a | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_organisation"></a> [organisation](#output\_organisation) | n/a |
<!-- END_TF_DOCS -->
