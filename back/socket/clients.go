package socket

import (
	"sync"

	"github.com/bettercallmolly/molly/socket/packets"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Mutex *sync.Mutex
	UUID  string
}

type Clients map[string]Client

var (
	ConnectedClients Clients
	MapMutex         sync.Mutex
)

func NewClients() Clients {
	return make(Clients)
}

func (c Clients) Add(conn *websocket.Conn, uuid string) {
	MapMutex.Lock()
	c[uuid] = Client{
		Conn:  conn,
		Mutex: &sync.Mutex{},
		UUID:  uuid,
	}
	MapMutex.Unlock()
	c.BroadcastExcept(uuid, packets.WriteWelcomePacket(uuid))
}

func (c Clients) Remove(uuid string) {
	MapMutex.Lock()
	delete(c, uuid)
	MapMutex.Unlock()
	c.BroadcastExcept(uuid, packets.WriteCyaPacket(uuid))
}

// Broadcast sends a message to all clients
func (c Clients) Broadcast(data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		client.Mutex.Lock()
		client.Conn.WriteMessage(websocket.BinaryMessage, data)
		client.Mutex.Unlock()
	}
	MapMutex.Unlock()
}

// BroadcastExcept sends a message to all clients except the one with the given uuid
func (c Clients) BroadcastExcept(uuid string, data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		if client.UUID != uuid {
			client.Mutex.Lock()
			client.Conn.WriteMessage(websocket.BinaryMessage, data)
			client.Mutex.Unlock()
		}
	}
	MapMutex.Unlock()
}
