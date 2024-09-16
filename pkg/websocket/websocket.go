package websocket

import (
	"encoding/json"
	"github.com/KianBahrami/PacketPirate/pkg/packets"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// websocket handle that accepts only connections from local computer
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://192.168.0.14:8000" || origin == "http://localhost:8000"
	},
}

// called if a request t /ws is done
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("New WebSocket connection request")
	// establish websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	stopChan := make(chan struct{}) // if closed all channel capture goroutines started afterwards are direvtly returning

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		} else if messageType == websocket.TextMessage && string(p) == "stop" {
			log.Println("Stopping packet capture")
			close(stopChan)                // next packet capture goroutine will get the stop chan signal
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
			go packets.CapturePackets(conn, &wg, stopChan, msg.Interface, msg.Filteroptions) // goroutine synchronized by wg
		}
	}
}
