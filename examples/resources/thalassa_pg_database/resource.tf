# Create a VPC for the database cluster
resource "thalassa_vpc" "example" {
  name            = "example-vpc"
  description     = "Example VPC for database cluster"
  region          = "nl-01"
  cidrs           = ["10.0.0.0/16"]
}

# Create a subnet for the database cluster
resource "thalassa_subnet" "example" {
  name            = "example-subnet"
  description     = "Example subnet for database cluster"
  vpc_id          = thalassa_vpc.example.id
  cidr            = "10.0.1.0/24"
}

# Create a database cluster for the PostgreSQL database
resource "thalassa_db_cluster" "example" {
  name                   = "example-db-cluster"
  description            = "Example database cluster for PostgreSQL database"
  subnet_id              = thalassa_subnet.example.id
  database_instance_type = "db-pgp-small"  # Available: db-pgp-small, db-pgp-medium, db-pgp-large, db-pgp-xlarge, db-pgp-2xlarge, db-pgp-4xlarge, db-dgp-small, db-dgp-medium, db-dgp-large, db-dgp-xlarge, db-dgp-2xlarge, db-dgp-4xlarge
  engine                 = "postgres"
  engine_version         = "15.13"
  allocated_storage      = 100
  volume_type_class      = "block"
}

# Create PostgreSQL roles first
resource "thalassa_pg_roles" "example" {
  db_cluster_id = thalassa_db_cluster.example.id
  name          = "myrole"
  password      = "secure_password_123" # Replace with secure password
}

# Create a PostgreSQL database with Thalassa default values
resource "thalassa_pg_database" "example" {
  name           = "mydatabase2"
  db_cluster_id  = thalassa_db_cluster.example.id
  owner_role_id  = thalassa_pg_roles.example.id
}

# Output the PostgreSQL database details
output "pg_database_id" {
  value = thalassa_pg_database.example.id
}

output "pg_database_name" {
  value = thalassa_pg_database.example.name
}
