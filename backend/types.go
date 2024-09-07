package main

import (
	"github.com/google/gopacket/pcap"
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
	Command   string `json:"command"`
	Interface string `json:"interface"`
}
