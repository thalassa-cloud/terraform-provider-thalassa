# Create VPCs for peering
resource "thalassa_vpc" "requester_vpc" {
  name   = "requester-vpc"
  region = "nl-01"
  cidrs  = ["10.0.0.0/16"]
}

resource "thalassa_vpc" "accepter_vpc" {
  name   = "accepter-vpc"
  region = "nl-01"
  cidrs  = ["10.1.0.0/16"]
}

# Create VPC peering connection
resource "thalassa_vpc_peering_connection" "example" {
  name             = "peering-connection-example"
  description      = "Peering connection between two VPCs in the same region"
  requester_vpc_id = thalassa_vpc.requester_vpc.id
  accepter_vpc_id  = thalassa_vpc.accepter_vpc.id
  auto_accept      = false
}

# # Example 1: Accept a peering connection by ID
resource "thalassa_vpc_peering_connection_acceptance" "accept_by_id" {
  peering_connection_id   = thalassa_vpc_peering_connection.example.id
  wait_for_active         = true
  wait_for_active_timeout = 1
}

# Configure route tables
data "thalassa_vpc_default_route_table" "requester_vpc" {
  vpc_id = thalassa_vpc.requester_vpc.id
}

data "thalassa_vpc_default_route_table" "accepter_vpc" {
  vpc_id = thalassa_vpc.accepter_vpc.id
}

# Configure route table routes
# Important; routes can only be created after the peering connection has been accepted and has become active
resource "thalassa_route_table_route" "requester_vpc" {
  route_table_id                = data.thalassa_vpc_default_route_table.requester_vpc.id
  destination_cidr              = thalassa_vpc.accepter_vpc.cidrs[0]
  target_vpc_peering_connection = thalassa_vpc_peering_connection_acceptance.accept_by_id.peering_connection_id
}

resource "thalassa_route_table_route" "accepter_vpc" {
  route_table_id                = data.thalassa_vpc_default_route_table.accepter_vpc.id
  destination_cidr              = thalassa_vpc.requester_vpc.cidrs[0]
  target_vpc_peering_connection = thalassa_vpc_peering_connection_acceptance.accept_by_id.peering_connection_id
}
