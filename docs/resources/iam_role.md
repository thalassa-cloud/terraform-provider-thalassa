---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "thalassa_iam_role Resource - terraform-provider-thalassa"
subcategory: ""
description: |-
  Manage an organisation role in Thalassa Cloud
---

# thalassa_iam_role (Resource)

Manage an organisation role in Thalassa Cloud



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the organisation role

### Optional

- `annotations` (Map of String) Annotations for the organisation role
- `description` (String) Description of the organisation role
- `labels` (Map of String) Labels for the organisation role

### Read-Only

- `created_at` (String) Creation timestamp of the organisation role
- `id` (String) The ID of this resource.
- `role_is_read_only` (Boolean) Whether the role is read-only and cannot be modified.
- `slug` (String) Slug of the organisation role
- `system` (Boolean) Whether the role is a system role
- `updated_at` (String) Last update timestamp of the organisation role
