package layers

import (
	"fmt"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

func tcpFlagsToString(tcp *layers.TCP) string {
	var flags []string
	if tcp.FIN {
		flags = append(flags, "FIN")
	}
	if tcp.SYN {
		flags = append(flags, "SYN")
	}
	if tcp.RST {
		flags = append(flags, "RST")
	}
	if tcp.PSH {
		flags = append(flags, "PSH")
	}
	if tcp.ACK {
		flags = append(flags, "ACK")
	}
	if tcp.URG {
		flags = append(flags, "URG")
	}
	if tcp.ECE {
		flags = append(flags, "ECE")
	}
	if tcp.CWR {
		flags = append(flags, "CWR")
	}
	if tcp.NS {
		flags = append(flags, "NS")
	}
	return fmt.Sprintf("%v", flags)
}

func ExtractPacketInfo(packet gopacket.Packet) types.PacketInfo {
	info := types.PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Length:    packet.Metadata().Length,
	}

	// Link Layer
	if linkLayer := packet.LinkLayer(); linkLayer != nil {
		info.LinkLayer.Protocol = linkLayer.LayerType().String()
		if eth, ok := linkLayer.(*layers.Ethernet); ok {
			info.LinkLayer.SrcMAC = eth.SrcMAC.String()
			info.LinkLayer.DstMAC = eth.DstMAC.String()
		}
	} 

	// maybe loop layer
	if loopbackLayer := packet.Layer(layers.LayerTypeLoopback); loopbackLayer != nil {
		info.LinkLayer.Protocol = loopbackLayer.LayerType().String()
		info.LinkLayer.SrcMAC = "-"
		info.LinkLayer.DstMAC = "-"
	}

	// maybe arp layer
	if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
		arp, _ := arpLayer.(*layers.ARP)
		info.ARPLayer.Operation = arp.Operation
		info.ARPLayer.SrcMAC = net.HardwareAddr(arp.SourceHwAddress).String()
		info.ARPLayer.DstMAC = net.HardwareAddr(arp.DstHwAddress).String()
		info.ARPLayer.SrcIP = net.IP(arp.SourceProtAddress).String()
		info.ARPLayer.DstIP = net.IP(arp.DstProtAddress).String()
	} else {
		info.ARPLayer.Operation = 42
		info.ARPLayer.SrcMAC = "None"
		info.ARPLayer.DstMAC = "None"
		info.ARPLayer.SrcIP = "None"
		info.ARPLayer.DstIP = "None"
	}

	// Network Layer
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		info.NetworkLayer.Protocol = networkLayer.LayerType().String()
		switch nl := networkLayer.(type) {
		case *layers.IPv4:
			info.NetworkLayer.SrcIP = nl.SrcIP.String()
			info.NetworkLayer.DstIP = nl.DstIP.String()
			info.NetworkLayer.TTL = nl.TTL
		case *layers.IPv6:
			info.NetworkLayer.SrcIP = nl.SrcIP.String()
			info.NetworkLayer.DstIP = nl.DstIP.String()
			info.NetworkLayer.TTL = nl.HopLimit
		}
	}

	// Transport Layer
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		info.TransportLayer.Protocol = transportLayer.LayerType().String()
		switch t := transportLayer.(type) {
		case *layers.TCP:
			info.TransportLayer.SrcPort = uint16(t.SrcPort)
			info.TransportLayer.DstPort = uint16(t.DstPort)
			info.TransportLayer.TCPFlags = tcpFlagsToString(t)
			info.TransportLayer.TCPSeq = t.Seq
			info.TransportLayer.TCPAck = t.Ack
			info.TransportLayer.TCPWindow = t.Window
		case *layers.UDP:
			info.TransportLayer.SrcPort = uint16(t.SrcPort)
			info.TransportLayer.DstPort = uint16(t.DstPort)
			info.TransportLayer.TCPFlags = "-"
			info.TransportLayer.TCPSeq = 42
			info.TransportLayer.TCPAck = 42
			info.TransportLayer.TCPWindow = 42
		}
	}

	// Application Layer
	if appLayer := packet.ApplicationLayer(); appLayer != nil {
		info.ApplicationLayer.Protocol = appLayer.LayerType().String()
		info.ApplicationLayer.PayloadSize = len(appLayer.Payload())
		info.ApplicationLayer.HTTPMethod = "-"
		info.ApplicationLayer.HTTPURL = "-"
		info.ApplicationLayer.HTTPVersion = "-"
	}

	info.Raw = packet.String()
	return info
}
