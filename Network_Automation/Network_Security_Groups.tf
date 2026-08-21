# Define local variables for modular architecture
locals {
  network_zones = {
    management = "10.100.1.0/24"
    production = "10.100.2.0/24"
    dmz        = "10.100.3.0/24"
  }
}

# Dynamically generate Subnets based on our local maps
resource "azurerm_subnet" "subnets" {
  for_each             = local.network_zones
  name                 = "subnet-${each.key}"
  resource_group_name  = azurerm_resource_group.net_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = [each.value]
}

# Create a Network Security Group (Access Control List)
resource "azurerm_network_security_group" "nsg" {
  name                = "nsg-prod-dmz-rules"
  location            = azurerm_resource_group.net_rg.location
  resource_group_name = azurerm_resource_group.net_rg.name

  # Inbound Security Rule allowing safe SSH management access
  security_rule {
    name                       = "Allow-SSH-Inbound"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "192.168.5.0/24" # Simulating your secure On-Prem Admin subnet
    destination_address_prefix = "*"
  }
}

# Bind the Security Group rules strictly to the DMZ Subnet
resource "azurerm_subnet_network_security_group_association" "association" {
  subnet_id                 = azurerm_subnet.subnets["dmz"].id
  network_security_group_id = azurerm_network_security_group.nsg.id
}
