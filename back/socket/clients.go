package socket

import (
	"math/rand"
	"sync"

	"github.com/bettercallmolly/molly/socket/packets"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Mutex    *sync.Mutex
	SocketId uint16
}

type Clients map[uint16]Client

var (
	ConnectedClients Clients
	MapMutex         sync.Mutex
)

func GenerateClientId() uint16 {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	var id uint16
	for {
		id = uint16(rand.Uint32())
		if _, ok := ConnectedClients[id]; !ok {
			return id
		}
	}
}

func NewClients() Clients {
	return make(Clients)
}

func (c Clients) Add(conn *websocket.Conn, socketId uint16) {
	MapMutex.Lock()
	c[socketId] = Client{
		Conn:     conn,
		Mutex:    &sync.Mutex{},
		SocketId: socketId,
	}
	MapMutex.Unlock()
	c.BroadcastExcept(socketId, packets.WriteWelcomePacket(socketId))
}

func (c Clients) Remove(socketId uint16) {
	MapMutex.Lock()
	delete(c, socketId)
	MapMutex.Unlock()
	c.BroadcastExcept(socketId, packets.WriteCyaPacket(socketId))
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
func (c Clients) BroadcastExcept(socketId uint16, data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		if client.SocketId != socketId {
			client.Mutex.Lock()
			client.Conn.WriteMessage(websocket.BinaryMessage, data)
			client.Mutex.Unlock()
		}
	}
	MapMutex.Unlock()
}
