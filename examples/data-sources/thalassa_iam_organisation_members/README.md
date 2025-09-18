# IAM Organisation Members Data Source Example

This example demonstrates how to use the `thalassa_iam_organisation_members` data source to retrieve information about all members of an organisation.

## Usage

The data source will return a list of organisation members with their details including:

- Member identity
- Creation timestamp
- Member type (OWNER or MEMBER)
- User information (subject, name, email, created_at)

### Filtering by Email

You can filter members by email address using the `email_filter` parameter:

```hcl
data "thalassa_iam_organisation_members" "filtered" {
  email_filter = "admin@example.com"
}
```

This will return only members whose email address matches the filter.

## Outputs

The example includes several outputs:

- `organisation_members`: Complete list of all organisation members
- `member_count`: Total number of members in the organisation
- `owners`: List of members with OWNER role
- `members`: List of members with MEMBER role
- `user_emails`: List of all user email addresses
- `filtered_members`: Members filtered by email address
- `filtered_member_count`: Number of members matching the email filter

## Running the Example

1. Configure your provider with appropriate credentials
2. Run `terraform init` to initialize the provider
3. Run `terraform plan` to see what will be created
4. Run `terraform apply` to execute the plan

## Notes

- The `organisation_id` parameter is optional. If not provided, the organisation from the provider configuration will be used.
- This data source is read-only and will not create or modify any resources.
