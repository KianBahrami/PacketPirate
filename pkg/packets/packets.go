package packets

import (
	"encoding/json"
	"fmt"
	"github.com/KianBahrami/PacketPirate/pkg/layers"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// sends json with packet info to the frontend via the websocket
func CapturePackets(conn *websocket.Conn, wg *sync.WaitGroup, stopChan chan struct{}, interfaceName string, filterOptions types.FilterOptions) {
	// new packet capture can be started if no other runs
	defer wg.Done()

	log.Printf("Attempting to open interface: %s", interfaceName)

	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Println("Error opening interface:", err)
		return
	}
	defer handle.Close()

	// Construct BPF filter string
	bpfFilter := ""
	if filterOptions.NetworkLayerProtocol != "any" && filterOptions.TransportLayerProtocol != "any" && filterOptions.NetworkLayerProtocol != "" && filterOptions.TransportLayerProtocol != "" {
		bpfFilter += fmt.Sprintf("%s and %s", filterOptions.NetworkLayerProtocol, filterOptions.TransportLayerProtocol)
	} else if filterOptions.TransportLayerProtocol != "any" {
		bpfFilter += fmt.Sprintf(filterOptions.TransportLayerProtocol)
	} else if filterOptions.NetworkLayerProtocol != "any" {
		bpfFilter += fmt.Sprintf(filterOptions.NetworkLayerProtocol)
	}
	if filterOptions.SrcIp != "" {
		// filter is not empty so we need an "and"
		if len(bpfFilter) != 0 {
			bpfFilter += fmt.Sprintf(" and src host %s", filterOptions.SrcIp)
		} else {
			bpfFilter += fmt.Sprintf("src host %s", filterOptions.SrcIp)
		}
	}
	if filterOptions.DestIp != "" {
		// filter is not empty so we need an "and"
		if len(bpfFilter) != 0 {
			bpfFilter += fmt.Sprintf(" and dst host %s", filterOptions.DestIp)
		} else {
			bpfFilter += fmt.Sprintf("host %s", filterOptions.DestIp)
		}
	}
	if filterOptions.MinPaylaodSize != "" {
		if len(bpfFilter) != 0 {
			bpfFilter += fmt.Sprintf(" and len > %v", filterOptions.MinPaylaodSize)
		} else {
			bpfFilter += fmt.Sprintf("len > %v", filterOptions.MinPaylaodSize)
		}

	}
	log.Printf("Applying BPF filter: %s", bpfFilter)

	// set filter
	err = handle.SetBPFFilter(bpfFilter)
	if err != nil {
		log.Println("Error setting BPF filter:", err)
		return
	}

	var bytesInAndOut int
	lastUpdateTime := time.Now()
	updateInterval := time.Second

	log.Println("Starting packet capture...")
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-stopChan:
			log.Println("Stopping packet capture...")
			return
		case packet := <-packetSource.Packets():
			log.Println("Packet captured")
			packetData := layers.ExtractPacketInfo(packet)
			
			// send packet to frontend
			jsonData, err := json.Marshal(packetData)
			if err != nil {
				log.Println("Error marshalling packet data:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Error sending packet data:", err)
				return
			}

			// compute BPS
			packetLength := packet.Metadata().Length
			bytesInAndOut += packetLength

			if time.Since(lastUpdateTime) >= updateInterval {
				bpsData := types.BPSInfo{
					Timestamp: time.Now().Unix(),
					BPS: float64(bytesInAndOut) / updateInterval.Seconds(),
				}
				bpsJSON, _ := json.Marshal(bpsData)
				if err := conn.WriteMessage(websocket.TextMessage, bpsJSON); err != nil {
					log.Println("Error sending datarate data:", err)
					return
				}
				bytesInAndOut = 0
				lastUpdateTime = time.Now()
			}
			

		}
	}
}
