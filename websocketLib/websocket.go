package websocketLib

import "log"

type WebsocketServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (w *WebsocketServer) Run() {
	for {
		select {
		case client := <-w.register:
			log.Printf("Registered with %v \n", client)
			w.registerClient(client)
		case client := <-w.unregister:
			log.Printf("Unregistered with %v \n", client)
			w.unregisterClient(client)
		case message := <-w.broadcast:
			log.Println("Mesah geldi gidiyor")
			w.broadcastToClients(message)
		}
	}
}

func (w *WebsocketServer) registerClient(c *Client) {
	w.clients[c] = true
}

func (w *WebsocketServer) unregisterClient(c *Client) {
	if _, ok := w.clients[c]; ok {
		delete(w.clients, c)
	}
}

func (w *WebsocketServer) broadcastToClients(message []byte) {
	for client := range w.clients {
		client.send <- message
	}
}

func (w *WebsocketServer) RegisterNewUser(client *Client) {
	w.register <- client
}
