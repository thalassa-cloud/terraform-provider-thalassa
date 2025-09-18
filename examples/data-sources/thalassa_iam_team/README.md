# IAM Team Data Source Example

This example demonstrates how to use the `thalassa_iam_team` data source to retrieve information about a team.

## Usage

The data source will return team information including:

- Team identity, name, slug, and description
- Labels and annotations
- Creation and update timestamps
- List of team members with their roles and user information

## Outputs

The example includes several outputs:

- `team_details`: Complete team information
- `team_members`: List of all team members
- `team_member_count`: Total number of team members
- `team_admin_members`: List of members with admin role

## Running the Example

1. Configure your provider with appropriate credentials
2. Run `terraform init` to initialize the provider
3. Run `terraform plan` to see what will be created
4. Run `terraform apply` to execute the plan

## Notes

- The `organisation_id` parameter is optional. If not provided, the organisation from the provider configuration will be used.
- You can search for a team by either `name` or `slug`.
- This data source is read-only and will not create or modify any resources.
