package interfaces

import (
	"encoding/json"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
)

// sends json with all available interfaces to frontend via http
func SendInterfaces(w http.ResponseWriter, r *http.Request) {
	// first get available interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Println("Error finding network interfaces:", err)
		return
	}

	log.Println("Available network interfaces:")
	for _, device := range devices {
		log.Printf("Name: %s, Description: %s", device.Name, device.Description)
	}

	// create struct that carries the interfaces
	availableifaces := types.AvailableInterfaces{}
	for _, device := range devices {
		iface := types.InterfaceData{
			Name:        device.Name,
			Addr:        device.Addresses,
			Description: device.Description,
			Flags:       device.Flags,
		}
		availableifaces.Interfaces = append(availableifaces.Interfaces, iface)
	}

	// set response content type to json
	w.Header().Set("Content-Type", "application/json")

	// send
	json.NewEncoder(w).Encode(availableifaces)
}
