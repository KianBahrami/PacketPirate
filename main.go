package main

import (
	"github.com/KianBahrami/PacketPirate/pkg/interfaces"
	"github.com/KianBahrami/PacketPirate/pkg/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
	"sync"
)

func main() {

	// setup two servers: one for backend, one for frontend
	var wg sync.WaitGroup
	wg.Add(2)

	// setup backend
	go func() {
		defer wg.Done()

		// setup routs and their functions
		http.HandleFunc("/ws", websocket.HandleWebSocket)
		http.HandleFunc("/get-interfaces", interfaces.SendInterfaces)

		// create CORS handler
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:8000"},		// "http://192.168.0.14:8000"
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

	// wait until both goroutines are finished
	wg.Wait()
}
