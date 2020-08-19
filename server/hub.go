package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("Register WS Client: %s", client.conn.RemoteAddr())
			h.clients[client] = true
		case client := <-h.unregister:
			log.Printf("Unregister WS Client: %s", client.conn.RemoteAddr())
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Printf("Broadcast to client %s: %s", client.conn.RemoteAddr(), message)
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}