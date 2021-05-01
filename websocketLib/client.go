package websocketLib

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	conn     *websocket.Conn
	wsServer *WebsocketServer
	send     chan []byte
	UserID   string `json:"user_id"`
}

func NewClient(conn *websocket.Conn, wsServer *WebsocketServer, id string) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		UserID:   id,
	}
}

func (c *Client) SendTest() {
	log.Println("Test Ediliyor LOG")
	c.send <- []byte("Test Ediliyor")
}

// Server'dan Client'a gelen
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			log.Println("c.sende girdi")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("Pong Ticker for %s \n", c.UserID)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Client'tan server'a giden
func (c *Client) ReadPump() {
	defer func() {
		c.wsServer.unregister <- c
		close(c.send)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetWriteDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			return
		}
		// client.wsServer.broadcast <- jsonMessage
		c.handleNewMessage(jsonMessage)
	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message messageSchema
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}
	switch message.TypeId {
	case UP:
		fmt.Printf("%s : Yukari Çıktı \n", c.UserID)
	case DOWN:
		fmt.Printf("%s : ,Aşağı İndi \n", c.UserID)
	case LEFT:
		fmt.Printf("%s : Sola Gitti \n", c.UserID)
	case RIGHT:
		fmt.Printf("%s : Sağa Gitti \n", c.UserID)
	}
	c.wsServer.broadcast <- message.Encode()
}
