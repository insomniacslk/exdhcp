package main

import (
	"flag"
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/netboot"
)

var (
	ver    = flag.Int("v", 6, "IP version to use")
	ifname = flag.String("i", "eth0", "Interface name")
	debug  = flag.Bool("d", false, "Print debug output")
)

func dhclient6(ifname string, verbose bool) error {
	llAddr, err := dhcpv6.GetLinkLocalAddr(ifname)
	if err != nil {
		return err
	}
	laddr := net.UDPAddr{
		IP:   llAddr,
		Port: dhcpv6.DefaultClientPort,
		Zone: ifname,
	}
	raddr := net.UDPAddr{
		IP:   dhcpv6.AllDHCPRelayAgentsAndServers,
		Port: dhcpv6.DefaultServerPort,
		Zone: ifname,
	}
	c := dhcpv6.NewClient()
	c.LocalAddr = &laddr
	c.RemoteAddr = &raddr
	conv, err := c.Exchange(ifname, nil)
	if verbose {
		for _, m := range conv {
			log.Print(m.Summary())
		}
	}
	if err != nil {
		return err
	}
	// extract the network configuration
	netconf, _, err := netboot.ConversationToNetconf(conv)
	if err != nil {
		return err
	}
	// configure the interface
	return netboot.ConfigureInterface(ifname, netconf)
}

func dhclient4(ifname string, verbose bool) error {
	client := dhcpv4.NewClient()
	conv, err := client.Exchange(ifname, nil)
	if verbose {
		for _, m := range conv {
			log.Print(m.Summary())
		}
	}
	if err != nil {
		return err
	}
	// extract the network configuration
	netconf, _, err := netboot.ConversationToNetconfv4(conv)
	if err != nil {
		return err
	}
	// configure the interface
	return netboot.ConfigureInterface(ifname, netconf)
}

func main() {
	flag.Parse()

	var err error
	if *ver == 6 {
		err = dhclient6(*ifname, *debug)
	} else {
		err = dhclient4(*ifname, *debug)
	}
	if err != nil {
		log.Fatal(err)
	}
}
