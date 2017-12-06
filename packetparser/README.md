# packetparser

An example of the encoding/decoding capabilities of the `dhcp` library, for both
DHCPv4 and DHCPv6. It can read a `pcap` file, or work as a client (which is the
default). See `./packetparser -h` for details, shown below for convenience:

```
$ ./packetparser -h
Usage of ./packetparser:
  -debug
        Enable debug output (default: false)
  -etherip
        Enables LayerTypeEtherIP instead of LayerTypeEthernet, use with linux-cooked PCAP files. (default: false)
  -i string
        Network interface to send packets through (default "eth0")
  -r string
        PCAP file to read from. If not specified, try to send an actual DHCP request
  -v int
        IP version to use (default 6)
```
