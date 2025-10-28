
# vpc
resource "thalassa_vpc" "vpc" {
  name   = "firewall-example"
  region = "nl-01"
  cidrs  = ["10.0.0.0/24"]
}

# Create a VPC firewall rule to allow SSH
resource "thalassa_vpc_firewall_rule" "allow_ssh" {
  vpc_id = thalassa_vpc.vpc.id
  name         = "allow-ssh"
  
  protocols {
    tcp = true
  }
  
  source      = "0.0.0.0/0"
  destination = "10.0.0.0/8"
  
  destination_ports = [22]
  action            = "allow"
  priority          = 100
  direction         = "inbound"
  state             = "active"
}

# Create a VPC firewall rule to allow HTTP/HTTPS
resource "thalassa_vpc_firewall_rule" "allow_web" {
  vpc_id = thalassa_vpc.vpc.id
  name         = "allow-web"
  
  protocols {
    tcp = true
  }
  
  source      = "0.0.0.0/0"
  destination = "10.0.0.0/8"
  
  destination_ports = [80, 443]
  action            = "allow"
  priority          = 200
  direction         = "inbound"
  state             = "active"
}
