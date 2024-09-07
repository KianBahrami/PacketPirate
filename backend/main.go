package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
    "github.com/KianBahrami/PacketPirate/pkg/layers"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

// websocket handle that accepts only connections from local computer
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://192.168.0.14:8000"
	},
}

func main() {
	// setup routs and their functions
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/get-interfaces", sendInterfaces)

	// create CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://192.168.0.14:8000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})
	handler := c.Handler(http.DefaultServeMux)

	// serve server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
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
		}

		// unpack message received via websocket
		var msg StartWebSocketMsg
		if err := json.Unmarshal(p, &msg); err != nil { // parse the message (p) into the struct StartWebSocketMsg
			log.Println("Error parsing websocket message:", err)
			continue
		}

		log.Printf("Received message: %s", string(p))

		if msg.Command == "start" { // received start message via websocket
			log.Println("Starting packet capture on interface:", msg.Interface)
			wg.Add(1)
			go capturePackets(conn, &wg, stopChan, msg.Interface) // goroutine synchronized by wg
		} else if messageType == websocket.TextMessage && string(p) == "stop" { // receive stop message via websocket
			log.Println("Stopping packet capture")
			close(stopChan)
			wg.Wait()                      // wait until every other capturePackets is finished
			stopChan = make(chan struct{}) // Reset the stop channel for the next capture session
		}
	}
}

// sends json with packet info to the frontend via the websocket
func capturePackets(conn *websocket.Conn, wg *sync.WaitGroup, stopChan chan struct{}, interfaceName string) {
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

	err = handle.SetBPFFilter("tcp")
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
	availableifaces := AvailableInterfaces{}
	for _, device := range devices {
		iface := InterfaceData{
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
