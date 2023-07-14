package socket

import (
	"github.com/bettercallmolly/molly/socket/packets"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn *websocket.Conn
	UUID string
}

type Clients map[string]*websocket.Conn

var (
	ConnectedClients Clients
)

func NewClients() Clients {
	return make(Clients)
}

func (c Clients) Add(conn *websocket.Conn, uuid string) {
	c[uuid] = conn
	c.BroadcastExcept(uuid, packets.WriteWelcomePacket(uuid))
}

func (c Clients) Remove(uuid string) {
	delete(c, uuid)
	c.BroadcastExcept(uuid, packets.WriteCyaPacket(uuid))
}

// Broadcast sends a message to all clients
func (c Clients) Broadcast(uuid string, data []byte) {
	for _, client := range c {
		if client != nil {
			client.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}

// BroadcastExcept sends a message to all clients except the one with the given uuid
func (c Clients) BroadcastExcept(uuid string, data []byte) {
	for id, client := range c {
		if client != nil && id != uuid {
			client.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
