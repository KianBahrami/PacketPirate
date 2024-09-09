package main

import (
	"encoding/json"
	"github.com/KianBahrami/PacketPirate/pkg/interfaces"
	"github.com/KianBahrami/PacketPirate/pkg/packets"
	"github.com/KianBahrami/PacketPirate/pkg/types"
	"log"
	"net/http"
	"sync"

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
		http.HandleFunc("/get-interfaces", interfaces.SendInterfaces)

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
			go packets.CapturePackets(conn, &wg, stopChan, msg.Interface, msg.Filteroptions) // goroutine synchronized by wg
		}
	}
}
