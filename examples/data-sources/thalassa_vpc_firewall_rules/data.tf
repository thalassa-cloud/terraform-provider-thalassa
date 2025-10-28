# Get all VPC firewall rules
data "thalassa_vpc_firewall_rules" "all_rules" {
  vpc_id = "vpc-12345678"
}

# Output all firewall rules
output "all_firewall_rules" {
  value = data.thalassa_vpc_firewall_rules.all_rules.firewall_rules
}

# Filter rules by action
locals {
  allow_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if rule.action == "allow"
  ]
  
  drop_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if rule.action == "drop"
  ]
  
  inbound_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if rule.direction == "inbound"
  ]
  
  outbound_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if rule.direction == "outbound"
  ]
  
  active_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if rule.state == "active"
  ]
}

# Output filtered rules
output "allow_rules" {
  value = local.allow_rules
}

output "drop_rules" {
  value = local.drop_rules
}

output "inbound_rules" {
  value = local.inbound_rules
}

output "outbound_rules" {
  value = local.outbound_rules
}

output "active_rules" {
  value = local.active_rules
}

# Count rules by different criteria
output "rule_counts" {
  value = {
    total_rules    = length(data.thalassa_vpc_firewall_rules.all_rules.firewall_rules)
    allow_rules    = length(local.allow_rules)
    drop_rules     = length(local.drop_rules)
    inbound_rules  = length(local.inbound_rules)
    outbound_rules = length(local.outbound_rules)
    active_rules   = length(local.active_rules)
  }
}

# Find rules with specific protocols
locals {
  tcp_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if length(rule.protocols) > 0 && rule.protocols[0].tcp == true
  ]
  
  udp_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if length(rule.protocols) > 0 && rule.protocols[0].udp == true
  ]
  
  icmp_rules = [
    for rule in data.thalassa_vpc_firewall_rules.all_rules.firewall_rules : rule
    if length(rule.protocols) > 0 && rule.protocols[0].icmp == true
  ]
}

output "protocol_rules" {
  value = {
    tcp_rules  = local.tcp_rules
    udp_rules  = local.udp_rules
    icmp_rules = local.icmp_rules
  }
}
