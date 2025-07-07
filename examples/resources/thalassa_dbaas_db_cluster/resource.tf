# Create a VPC for the database cluster
resource "thalassa_vpc" "example" {
  name        = "example-vpc"
  description = "Example VPC for database cluster"
  region      = "nl-01"
  cidrs       = ["10.0.0.0/16"]
}

# Create a subnet for the database cluster
resource "thalassa_subnet" "example" {
  name        = "example-subnet"
  description = "Example subnet for database cluster"
  vpc_id      = thalassa_vpc.example.id
  cidr        = "10.0.1.0/24"
}

# Create a security group for the DB cluster
resource "thalassa_security_group" "example" {
  name         = "example-db-security-group"
  description  = "Example security group for DB cluster"
  vpc_identity = thalassa_vpc.example.id
}

# Create a database cluster with Thalassa default values
resource "thalassa_dbaas_db_cluster" "example" {
  name                   = "example-db-cluster"
  description            = "Example database cluster for documentation"
  subnet_id              = thalassa_subnet.example.id
  database_instance_type = "db-pgp-small" # Available: db-pgp-small, db-pgp-medium, db-pgp-large, db-pgp-xlarge, db-pgp-2xlarge, db-pgp-4xlarge, db-dgp-small, db-dgp-medium, db-dgp-large, db-dgp-xlarge, db-dgp-2xlarge, db-dgp-4xlarge
  engine                 = "postgres"
  engine_version         = "15.13"
  allocated_storage      = 100
  volume_type_class      = "block"
}

# Output the database cluster details
output "db_cluster_id" {
  value = thalassa_dbaas_db_cluster.example.id
}

output "db_cluster_name" {
  value = thalassa_dbaas_db_cluster.example.name
}

output "db_cluster_endpoint" {
  value = thalassa_dbaas_db_cluster.example.endpoint_ipv4
}

output "db_cluster_port" {
  value = thalassa_dbaas_db_cluster.example.port
}
