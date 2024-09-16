package types

import (
	"github.com/google/gopacket/pcap"
	"time"
)

// struct for storing information about an interface
type InterfaceData struct {
	Name        string                  `json:"name"`
	Addr        []pcap.InterfaceAddress `json:"addr"`
	Description string                  `json:"description"`
	Flags       uint32                  `json:"flags"`
}

// struct to save the available interfaces
type AvailableInterfaces struct {
	Interfaces []InterfaceData `json:"interfaces"`
}

// message when starting the webservice
type StartWebSocketMsg struct {
	Command       string        `json:"command"`
	Interface     string        `json:"interface"`
	Filteroptions FilterOptions `json:"filteroptions"`
}

type FilterOptions struct {
	NetworkLayerProtocol   string `json:"networklayerprotocol"`
	TransportLayerProtocol string `json:"transportlayerprotocol"`
	SrcIp                  string `json:"srcip"`
	DestIp                 string `json:"destip"`
	MinPaylaodSize         string `json:"minpayloadsize"`
}

type BPSInfo struct {
	Timestamp int64   `json:"timestamp"`
	BPS       float64 `json:"bps"`
}

type PacketInfo struct {
	Timestamp time.Time `json:"time"`
	Length    int       `json:"length"`
	LinkLayer struct {
		Protocol string `json:"protocol"`
		SrcMAC   string `json:"src"`
		DstMAC   string `json:"dest"`
	} `json:"linklayer"`
	NetworkLayer struct {
		Protocol     string `json:"protocol"`
		SrcIP        string `json:"src"`
		DstIP        string `json:"dest"`
		TTL          uint8  `json:"ttl"`
		ARPOperation uint16 `json:"arpoperation,omitempty"`
	} `json:"networklayer"`
	ARPLayer struct {
		Operation uint16 `json:"operation,omitempty"`
		SrcMAC    string `json:"srcmac,omitempty"`
		DstMAC    string `json:"dstmac,omitempty"`
		SrcIP     string `json:"srcip,omitempty"`
		DstIP     string `json:"dstip,omitempty"`
	} `json:"arplayer,omitempty"`
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