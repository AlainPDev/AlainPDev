# Configure the Microsoft Azure Provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# Create a logical Resource Group container
resource "azurerm_resource_group" "net_rg" {
  name     = "rg-enterprise-networking-prod"
  location = "East US"
}

# Provision a Core Virtual Network
resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-core-prod-001"
  address_space       = ["10.100.0.0/16"]
  location            = azurerm_resource_group.net_rg.location
  resource_group_name = azurerm_resource_group.net_rg.name
}
