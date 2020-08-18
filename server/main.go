package main

import (
	sdk "agones.dev/agones/sdks/go"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	go doSignal()
	port := flag.String("port", "7654", "The port to listen to http traffic on")

	log.Print("Creating SDK instance")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}

	log.Print("Starting Health Ping")
	stop := make(chan struct{})
	go doHealth(s, stop)

	log.Print("Starting WS Hub")
	hub := newHub()
	go hub.run()
	ready(s)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		go client.writePump()
		go client.readPump()
	})

	log.Print("Starting HTTP Server")
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
