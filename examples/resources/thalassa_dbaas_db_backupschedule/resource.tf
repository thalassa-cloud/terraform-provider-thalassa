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

# Create a database cluster for the backup schedule
resource "thalassa_dbaas_db_cluster" "example" {
  name                   = "example-db-cluster"
  description            = "Example database cluster for backup schedule"
  subnet_id              = thalassa_subnet.example.id
  database_instance_type = "db-pgp-small" # Available: db-pgp-small, db-pgp-medium, db-pgp-large, db-pgp-xlarge, db-pgp-2xlarge, db-pgp-4xlarge, db-dgp-small, db-dgp-medium, db-dgp-large, db-dgp-xlarge, db-dgp-2xlarge, db-dgp-4xlarge
  engine                 = "postgres"
  engine_version         = "15.13"
  allocated_storage      = 100
  volume_type_class      = "block"
}

# Create a database backup schedule with Thalassa default values
resource "thalassa_dbaas_db_backupschedule" "example" {
  db_cluster_id    = thalassa_dbaas_db_cluster.example.id
  name             = "example-backup-schedule"
  schedule         = "0 2 * * *" # Daily at 2 AM
  retention_policy = "7d"        # Available: 7d, 14d, 30d, 90d, 180d, 365d, 730d
}

# Output the backup schedule details
output "backup_schedule_id" {
  value = thalassa_dbaas_db_backupschedule.example.id
}

output "backup_schedule_name" {
  value = thalassa_dbaas_db_backupschedule.example.name
}
