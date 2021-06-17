package vnet

import "encoding/xml"

// VNet represents a virtual network.
type VNet struct {
	XMLName   xml.Name `xml:"VNET"`
	Name      string   `xml:"NAME"`
	ID        int      `xml:"ID"`
	Cluster   string   `xml:"CLUSTER"`
	ClusterID int      `xml:"CLUSTER_ID"`
	Bridge    string   `xml:"BRIDGE"`
}

// NIC represents a network interface.
type NIC struct {
	XMLName   xml.Name `xml:"NIC"`
	Name      string   `xml:"NETWORK"`
	NetworkID int      `xml:"NETWORK_ID"`
}
