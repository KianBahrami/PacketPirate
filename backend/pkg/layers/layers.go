package layers

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PacketInfo struct {
	Timestamp time.Time `json:"time"`
	Length    int       `json:"length"`
	LinkLayer struct {
		Protocol string `json:"protocol"`
		SrcMAC   string `json:"src"`
		DstMAC   string `json:"dest"`
	} `json:"linklayer"`
	NetworkLayer struct {
		Protocol string `json:"protocol"`
		SrcIP    string `json:"src"`
		DstIP    string `json:"dest"`
		TTL      uint8  `json:"ttl"`
	} `json:"networklayer"`
	TransportLayer struct {
		Protocol  string `json:"protocol"`
		SrcPort   uint16 `json:"src"`
		DstPort   uint16 `json:"dest"`
		TCPFlags  string `json:"tcpflags"`
		TCPSeq    uint32 `json:"tcpseq"`
		TCPAck    uint32 `json:"tcpack"`
		TCPWindow uint16 `json:"tcpwindow"`
	} `json:"transportlayer"`
	ApplicationLayer struct {
		Protocol    string `json:"protocol"`
		PayloadSize int    `json:"payloadsize"`
		HTTPMethod  string `json:"httpmethod"`
		HTTPURL     string `json:"httpurl"`
		HTTPVersion string `json:"httpversion"`
	} `json:"applicationlayer"`
	Raw string `json:"raw"`
}

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

func ExtractPacketInfo(packet gopacket.Packet) PacketInfo {
	info := PacketInfo{
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

	// Network Layer
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		info.NetworkLayer.Protocol = networkLayer.LayerType().String()
		if ip4, ok := networkLayer.(*layers.IPv4); ok {
			info.NetworkLayer.SrcIP = ip4.SrcIP.String()
			info.NetworkLayer.DstIP = ip4.DstIP.String()
			info.NetworkLayer.TTL = ip4.TTL
		} else if ip6, ok := networkLayer.(*layers.IPv6); ok {
			info.NetworkLayer.SrcIP = ip6.SrcIP.String()
			info.NetworkLayer.DstIP = ip6.DstIP.String()
			info.NetworkLayer.TTL = ip6.HopLimit
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
		// TODO
		//    // Example: Parse HTTP if present
		//    if httpLayer := packet.Layer(layers.LayerTypeHTTP); httpLayer != nil {
		//        if http, ok := httpLayer.(*layers.HTTP); ok {
		//            info.ApplicationLayer.HTTPMethod = string(http.Method)
		//            info.ApplicationLayer.HTTPURL = string(http.URL)
		//            info.ApplicationLayer.HTTPVersion = string(http.Version)
		//        }
		//    }
	}

	info.Raw = packet.String()
	return info
}
