package main

import (
	"encoding/json"
	"fmt"
	"github.com/KianBahrami/PacketPirate/pkg/layers"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"log"
	"net/http"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// websocket handle that accepts only connections from local computer
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://192.168.0.14:8000" || origin == "http://localhost:8000"
	},
}

func main() {

	// setup two servers: one for backend, one for frontend
	var wg sync.WaitGroup
	wg.Add(2)

	// setup backend
	go func() {
		defer wg.Done()

		// setup routs and their functions
		http.HandleFunc("/ws", handleWebSocket)
		http.HandleFunc("/get-interfaces", sendInterfaces)

		// create CORS handler
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://192.168.0.14:8000", "http://localhost:8000"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		})
		handler := c.Handler(http.DefaultServeMux)

		// serve server
		log.Println("Backend server starting on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", handler))
	}()

	// setup frontend
	go func() {
		defer wg.Done()

		// Serve html
		fs := http.FileServer(http.Dir("./static"))

		frontendMux := http.NewServeMux()
		frontendMux.Handle("/", fs)

		log.Println("Frontend server starting on http://localhost:8000")
		if err := http.ListenAndServe(":8000", frontendMux); err != nil {
			log.Fatal("Frontend server error:", err)
		}
	}()

	wg.Wait()
}

// called if a request t /ws is done
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("New WebSocket connection request")
	// establish websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		} else if messageType == websocket.TextMessage && string(p) == "stop" {
			log.Println("Stopping packet capture")
			close(stopChan)
			wg.Wait()                      // wait until every other capturePackets is finished
			stopChan = make(chan struct{}) // Reset the stop channel for the next capture session
		}

		// unpack message received via websocket
		var msg types.StartWebSocketMsg
		if err := json.Unmarshal(p, &msg); err != nil { // parse the message (p) into the struct StartWebSocketMsg
			log.Println("Error parsing websocket message:", err)
			continue
		}

		log.Printf("Received message: %s", string(p))

		if msg.Command == "start" { // received start message via websocket
			log.Println("Starting packet capture on interface:", msg.Interface)
			wg.Add(1)
			go capturePackets(conn, &wg, stopChan, msg.Interface, msg.Filteroptions) // goroutine synchronized by wg
		}
	}
}

// sends json with packet info to the frontend via the websocket
func capturePackets(conn *websocket.Conn, wg *sync.WaitGroup, stopChan chan struct{}, interfaceName string, filterOptions types.FilterOptions) {
	// new packet capture can be started if no other runs
	defer wg.Done()

	// You may need to change this to match an interface on your system
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
	log.Printf("Applying BPF filter: %s", bpfFilter)

	// set filter
	err = handle.SetBPFFilter(bpfFilter)
	if err != nil {
		log.Println("Error setting BPF filter:", err)
		return
	}

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

			jsonData, err := json.Marshal(packetData)
			if err != nil {
				log.Println("Error marshalling packet data:", err)
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Println("Error sending packet data:", err)
				return
			}
		}
	}
}

// sends json with all available interfaces to frontend via http
func sendInterfaces(w http.ResponseWriter, r *http.Request) {
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
