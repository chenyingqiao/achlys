package cmd

import (
	"net"

	"github.com/vishvananda/netlink"
)

type Network struct{
	Name string
	IpRange *net.IPNet
	Driver string
}

type Endpoint struct{
	ID string
	Device netlink.Veth
	IPAddress net.IP
	MacAddress net.HardwareAddr
	PortMapping []string
	Network *Network
}

type Driver interface{
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network)
	Connect(network *Network,endpoint *Endpoint)
	Disconnect(network Network, endpoint *Endpoint)
}

func CreateNetwork(driver, subnet, name string) error {
	// _, cidr, _ := net.ParseCIDR(subnet)
	return nil
}
