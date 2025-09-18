# IAM Team Resource Example

This example demonstrates how to create and manage teams using the Thalassa Cloud Terraform provider.

## Features

- Create teams with custom names and descriptions
- Add labels and annotations for better organization
- Manage team metadata and properties
- Add team members by user identity or email address
- Automatic email-to-user-identity resolution

## Usage

1. Configure your Thalassa Cloud provider credentials
2. Run the example:

```bash
terraform init
terraform plan
terraform apply
```

## Configuration

The example creates a team with the following properties:

- **Name**: `example-team`
- **Description**: `An example team for demonstration purposes`
- **Labels**: Environment and project labels for categorization
- **Annotations**: Contact and owner information
- **Members**: Team members added by both user identity and email address

## Outputs

The example provides the following outputs:

- `team_id`: The unique identifier of the created team
- `team_name`: The name of the team
- `team_slug`: The URL-friendly slug of the team
- `team_description`: The description of the team
- `team_members`: List of team members with their roles
- `team_member_count`: Total number of team members

## Cleanup

To clean up the resources created by this example:

```bash
terraform destroy
```

## Notes

- Teams are organization-scoped resources
- Team names must be unique within an organization
- Labels and annotations are optional but recommended for better resource management
- Team members can be specified by either `user_identity` or `email` (but not both)
- When using `email`, the user identity is automatically resolved from organisation members
- Member roles can be "OWNER", "ADMIN", or "MEMBER" (defaults to "MEMBER" if not specified)
- When only providing an email, the role defaults to "MEMBER" automatically 