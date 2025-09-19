# Terraform Provider Thalassa Cloud

Thalassa Cloud Terraform Provider

## Documentation

- [registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs](https://registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install)

## Examples

### Provider Configuration

- [Basic Provider Setup](./examples/provider/) - Example showing how to configure the Thalassa provider with authentication

### Infrastructure (IaaS)

#### Networking
- [VPC](./examples/resources/thalassa_vpc/) - Create a VPC with subnets
- [Subnet](./examples/resources/thalassa_subnet/) - Create subnets within a VPC
- [Security Group](./examples/resources/thalassa_security_group/) - Configure security groups
- [Route Table](./examples/resources/thalassa_route_table/) - Set up routing tables
- [Route Table Route](./examples/resources/thalassa_route_table_route/) - Configure routes within route tables
- [NAT Gateway](./examples/resources/thalassa_natgateway/) - Set up NAT gateways for outbound internet access

#### Compute
- [Virtual Machine Instance](./examples/resources/thalassa_virtual_machine_instance/) - Deploy virtual machines
- [Block Volume](./examples/resources/thalassa_block_volume/) - Create and manage block storage volumes
- [Block Volume Attachment](./examples/resources/thalassa_block_volume_attachment/) - Attach volumes to instances
- [Cloud Init Template](./examples/resources/thalassa_cloud_init_template/) - Create cloud-init templates for instance initialization

#### Load Balancing
- [Load Balancer](./examples/resources/thalassa_loadbalancer/) - Set up load balancers
- [Load Balancer Listener](./examples/resources/thalassa_loadbalancer_listener/) - Configure load balancer listeners
- [Target Group](./examples/resources/thalassa_target_group/) - Create target groups for load balancers
- [Target Group Attachment](./examples/resources/thalassa_target_group_attachment/) - Attach resources to target groups

#### Kubernetes
- [Kubernetes Cluster](./examples/resources/thalassa_kubernetes/) - Deploy and manage Kubernetes clusters with node pools

#### Database as a Service (DBaaS)
- [PostgreSQL Database](./examples/resources/thalassa_dbaas_pg_database/) - Create PostgreSQL databases
- [PostgreSQL Roles](./examples/resources/thalassa_dbaas_pg_roles/) - Manage PostgreSQL database roles
- [Database Cluster](./examples/resources/thalassa_dbaas_db_cluster/) - Set up database clusters
- [Database Backup Schedule](./examples/resources/thalassa_dbaas_db_backupschedule/) - Configure automated backup schedules

#### Identity & Access Management (IAM)
- [IAM Role](./examples/resources/thalassa_iam_role/) - Create and manage IAM roles
- [IAM Role Binding](./examples/resources/thalassa_iam_role_binding/) - Bind roles to service accounts or users
- [Service Account](./examples/resources/thalassa_iam_service_account/) - Create and manage IAM service accounts
- [Service Account Access Credential](./examples/resources/thalassa_iam_service_account_access_credential/) - Manage access credentials for service accounts

## License

[Apache 2.0 License 2.0](lICENSE)

## Contributing

Set up the provider locally and make sure to change the variables if needed:
```bash
make install NAMESPACE=local HOSTNAME=terraform.local OS_ARCH=darwin_arm64
```

Use the locally installed provider in your Terraform configuration:
```hcl
terraform {
  required_providers {
    thalassa = {
      source = "thalassa.cloud/thalassa/thalassa"
    }
  }
}
```
