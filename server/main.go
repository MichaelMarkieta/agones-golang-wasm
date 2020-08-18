package main

import (
	sdk "agones.dev/agones/sdks/go"
	"log"
	"net/http"
)

func main() {
	go doSignal()

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
		log.Print("WS request made")
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
	log.Fatal(http.ListenAndServe(":7654", nil))
}
