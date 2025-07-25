---
page_title: "thalassa_cloud_init_template Resource - terraform-provider-thalassa"
subcategory: "Compute"
description: |-
  
---

# thalassa_cloud_init_template (Resource)



## Example Usage

```terraform
# Create a cloud init template with Thalassa default values
resource "thalassa_cloud_init_template" "example" {
  name    = "example-cloud-init-template"
  content = "#cloud-config\npackage_update: true\npackage_upgrade: true\npackages:\n  - nginx\n  - curl"
}

# Output the cloud init template details
output "cloud_init_template_id" {
  value = thalassa_cloud_init_template.example.id
}

output "cloud_init_template_name" {
  value = thalassa_cloud_init_template.example.name
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) The content of the cloud init template
- `name` (String) The name of the cloud init template

### Optional

- `annotations` (Map of String) Annotations to add to the cloud init template
- `labels` (Map of String) Labels to add to the cloud init template
- `organisation_id` (String) Reference to the Organisation of the Machine Type. If not provided, the organisation of the (Terraform) provider will be used.

### Read-Only

- `id` (String) The identity of the cloud init template
- `slug` (String) The slug of the cloud init template

 