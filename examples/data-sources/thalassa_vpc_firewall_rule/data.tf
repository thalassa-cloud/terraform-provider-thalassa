# Get VPC firewall rule by name
data "thalassa_vpc_firewall_rule" "by_name" {
  vpc_id = "vpc-12345678"
  name   = "allow-ssh"
}

# Get VPC firewall rule by identity
data "thalassa_vpc_firewall_rule" "by_identity" {
  vpc_id = "vpc-12345678"
  identity = "firewall-rule-12345678"
}

# Output the firewall rule details
output "firewall_rule_by_name" {
  value = {
    id          = data.thalassa_vpc_firewall_rule.by_name.id
    name        = data.thalassa_vpc_firewall_rule.by_name.name
    action      = data.thalassa_vpc_firewall_rule.by_name.action
    priority    = data.thalassa_vpc_firewall_rule.by_name.priority
    direction   = data.thalassa_vpc_firewall_rule.by_name.direction
    state       = data.thalassa_vpc_firewall_rule.by_name.state
    protocols   = data.thalassa_vpc_firewall_rule.by_name.protocols
    source      = data.thalassa_vpc_firewall_rule.by_name.source
    destination = data.thalassa_vpc_firewall_rule.by_name.destination
  }
}

output "firewall_rule_by_identity" {
  value = {
    id          = data.thalassa_vpc_firewall_rule.by_identity.id
    name        = data.thalassa_vpc_firewall_rule.by_identity.name
    action      = data.thalassa_vpc_firewall_rule.by_identity.action
    priority    = data.thalassa_vpc_firewall_rule.by_identity.priority
    direction   = data.thalassa_vpc_firewall_rule.by_identity.direction
    state       = data.thalassa_vpc_firewall_rule.by_identity.state
    protocols   = data.thalassa_vpc_firewall_rule.by_identity.protocols
    source      = data.thalassa_vpc_firewall_rule.by_identity.source
    destination = data.thalassa_vpc_firewall_rule.by_identity.destination
  }
}
