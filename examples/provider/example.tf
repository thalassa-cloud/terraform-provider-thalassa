# terraform {
#   required_providers {
#     thalassa = {
#       version = ">= 0.1"
#       source  = "thalassa-cloud/thalassa"
#     }
#   }
# }

terraform {
  required_providers {
    thalassa = {
      source = "thalassa.cloud/thalassa/thalassa"
    }
  }
}

variable "thalassa_token" {
  sensitive = true
}

variable "thalassa_api" {
}

provider "thalassa" {
  token           = var.thalassa_token
  api             = var.thalassa_api
  organisation_id = "containerinfra-xxbvs"
}

data "thalassa_organisation" "this" {
  slug = "containerinfra-xxbvs"
}


output "organisation_id" {
  value = data.thalassa_organisation.this
}

# data "thalassa_regions" "this" {
#   organisation = data.thalassa_organisation.this.slug
# }

# output "regions" {
#   value = data.thalassa_regions.this
# }

# vpc
resource "thalassa_vpc" "this" {
  name            = "localdev"
  organisation_id = data.thalassa_organisation.this.id
  cidrs           = ["10.0.0.0/16", "10.2.0.0/16", "10.3.0.0/16"]
  region          = "nl-01"
}

# # subnet
resource "thalassa_subnet" "subnet" {
  name            = "localdev"
  organisation_id = data.thalassa_organisation.this.id
  vpc_id          = thalassa_vpc.this.id
  cidr            = "10.2.0.0/24"
}

resource "thalassa_natgateway" "nat_gateway" {
  name            = "localdev"
  organisation_id = data.thalassa_organisation.this.id
  subnet_id       = thalassa_subnet.subnet.id
}

resource "thalassa_route_table" "route_table" {
  name            = "localdev"
  organisation_id = data.thalassa_organisation.this.id
  vpc_id          = thalassa_vpc.this.id
}

resource "thalassa_route_table_route" "route" {
  route_table_id   = thalassa_route_table.route_table.id
  destination_cidr = "0.0.0.0/0"
  gateway_address  = thalassa_natgateway.nat_gateway.endpoint_ip
}

# resource "thalassa_kubernetes_cluster" "this" {
#   name         = "my-cluster2"
#   organisation_id = data.thalassa_organisation.this.id
#   subnet_id    = thalassa_subnet.subnet.id

#   cluster_version         = "v1.31.5-0"
#   networking_cni          = "cilium"
#   networking_pod_cidr     = "10.4.0.0/16"
#   networking_service_cidr = "10.5.0.0/16"
# }

# resource "thalassa_kubernetes_node_pool" "worker" {
#   name               = "worker"
#   organisation       = data.thalassa_organisation.this.id
#   cluster            = thalassa_kubernetes_cluster.this.id
#   replicas           = 1
#   machine_type       = "pgp-large"
#   subnet             = thalassa_subnet.subnet.id
#   # kubernetes_version = "v1.31.5-0"
#   upgrade_strategy   = "always"

#   enable_autohealing = true
# }

# data "thalassa_block_volume_type" "example" {
#   name         = "example"
#   organisation = data.thalassa_organisation.this.id
#   region       = "nl-01"
# }

# resource "thalassa_block_volume" "example" {
#   name            = "example2"
#   organisation_id = data.thalassa_organisation.this.id
#   region          = "nl-01"
#   size_gb         = 15
#   volume_type     = "10b8e6d0-bb52-40b4-be2c-012a05058064"
# }

data "thalassa_machine_image" "ubuntu_22_04_01" {
  name = "ubuntu-22-04-01"
}

resource "thalassa_virtual_machine_instance" "testvm" {
  name                = "testvm"
  organisation_id     = data.thalassa_organisation.this.id
  machine_type        = "pgp-xlarge"
  machine_image       = data.thalassa_machine_image.ubuntu_22_04_01.id
  subnet_id           = thalassa_subnet.subnet.id
  # root_volume_id      = thalassa_block_volume.example.id
  root_volume_type    = "10b8e6d0-bb52-40b4-be2c-012a05058064"
  availability_zone   = "nl-01a"
  root_volume_size_gb = 30
}
