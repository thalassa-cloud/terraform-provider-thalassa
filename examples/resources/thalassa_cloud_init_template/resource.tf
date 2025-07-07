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

