package socket

import (
	"math/rand"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Mutex    *sync.Mutex
	SocketId uint32
}

type Clients map[uint32]Client

var (
	NbClients        uint32
	ConnectedClients Clients
	MapMutex         sync.Mutex
)

func GenerateClientId() uint32 {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	var id uint32
	for {
		id = rand.Uint32()
		if _, ok := ConnectedClients[id]; !ok {
			return id
		}
	}
}

func NewClients() Clients {
	return make(Clients)
}

func (c Clients) Add(conn *websocket.Conn, socketId uint32) {
	MapMutex.Lock()
	c[socketId] = Client{
		Conn:     conn,
		Mutex:    &sync.Mutex{},
		SocketId: socketId,
	}
	NbClients++
	MapMutex.Unlock()
}

func (c Clients) Remove(socketId uint32) {
	MapMutex.Lock()
	NbClients--
	delete(c, socketId)
	MapMutex.Unlock()
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
func (c Clients) BroadcastExcept(socketId uint32, data []byte) {
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

func GetNumberOfClients() uint32 {
	return NbClients
}
