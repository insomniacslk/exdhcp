package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"io"
	"log"
	"net"
)

var ver = flag.Int("v", 6, "IP version to use")
var infile = flag.String("r", "", "PCAP file to read from. If not specified, try to send an actual DHCP request")
var iface = flag.String("i", "eth0", "Network interface to send packets through")
var useEtherIP = flag.Bool("etherip", false, "Enables LayerTypeEtherIP instead of LayerTypeEthernet, use with linux-cooked PCAP files. (default: false)")
var debug = flag.Bool("debug", false, "Enable debug output (default: false)")

func Clientv4() {
	client := dhcpv4.Client{}
	conv, err := client.Exchange(nil, "wlp58s0")
	// don't exit immediately if there's an error, since `conv` will always
	// contain at least the SOLICIT message. So print it out first
	for _, m := range conv {
		log.Print(m.Summary())
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Clientv6() {
	llAddr, err := dhcpv6.GetLinkLocalAddr(*iface)
	if err != nil {
		panic(err)
	}
	laddr := net.UDPAddr{
		IP:   *llAddr,
		Port: 546,
		Zone: *iface,
	}
	raddr := net.UDPAddr{
		IP:   dhcpv6.AllDHCPRelayAgentsAndServers,
		Port: 547,
		Zone: *iface,
	}
	c := dhcpv6.Client{
		LocalAddr:  &laddr,
		RemoteAddr: &raddr,
	}
	conv, err := c.Exchange("wlp58s0", nil)
	// don't exit immediately if there's an error, since `conv` will always
	// contain at least the SOLICIT message. So print it out first
	for _, m := range conv {
		fmt.Print(m.Summary())
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	if *infile == "" {
		if *ver == 4 {
			Clientv4()
		} else {
			Clientv6()
		}
	} else {
		handle, err := pcap.OpenOffline(*infile)
		if err != nil {
			panic(err)
		}
		defer handle.Close()
		var pcapFilter string
		if *ver == 6 {
			pcapFilter = "ip6 and udp portrange 546-547"
		} else {
			pcapFilter = "ip and udp portrange 67-68"
		}
		err = handle.SetBPFFilter(pcapFilter)
		if err != nil {
			panic(err)
		}
		var layerType gopacket.LayerType
		if *useEtherIP {
			layerType = layers.LayerTypeEtherIP
		} else {
			layerType = layers.LayerTypeEthernet
		}
		for {
			data, _, err := handle.ReadPacketData()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			pkt := gopacket.NewPacket(data, layerType, gopacket.Default)
			if *debug {
				fmt.Println(pkt)
			}
			if udpLayer := pkt.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				if *debug {
					fmt.Println(udp.Payload)
				}
				if *ver == 4 {
					d, err := dhcpv4.FromBytes(udp.Payload)
					if err != nil {
						panic(err)
					}
					fmt.Println(d.Summary())
				} else {
					d, err := dhcpv6.FromBytes(udp.Payload)
					if err != nil {
						panic(err)
					}
					fmt.Println(d.Summary())
				}
			}
		}
	}
}
