package socket

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var (
	LastSocketId uint32
)

const (
	S_CURSORS = 0x00

	MASK_CURSORS = 1 << S_CURSORS
)

type Subscriptions struct {
	Cursors bool
}

func (s *Subscriptions) ToBitmask() uint8 {
	var mask uint8
	if s.Cursors {
		mask |= MASK_CURSORS
	}
	return mask
}

type Client struct {
	Conn          *websocket.Conn
	Mutex         *sync.Mutex
	SocketId      uint32
	Subscriptions *Subscriptions
}

func (c *Client) Subscribe(channel uint8) error {
	switch channel {
	case S_CURSORS:
		c.Subscriptions.Cursors = true
	default:
		return ErrInvalidChannel
	}
	return nil
}

func (c *Client) Unsubscribe(channel uint8) error {
	switch channel {
	case S_CURSORS:
		c.Subscriptions.Cursors = false
	default:
		return ErrInvalidChannel
	}
	return nil
}

type Clients map[uint32]Client

var (
	NbClients        uint32
	ConnectedClients Clients
	MapMutex         sync.Mutex

	// Use to store clients that have been disconnected while a Broadcast or BroadcastExcept call
	deadClients []uint32
)

func GenerateClientId() uint32 {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	// Increment the last socket id until it is not used
	for {
		LastSocketId++
		if _, ok := ConnectedClients[LastSocketId]; !ok {
			return LastSocketId
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
		Subscriptions: &Subscriptions{
			Cursors: false,
		},
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

func (c Clients) RemoveLocked(socketId uint32) {
	NbClients--
	delete(c, socketId)
}

// Broadcast sends a message to all clients
func (c Clients) Broadcast(data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		client.Mutex.Lock()
		err := client.Conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			deadClients = append(deadClients, client.SocketId)
		}
		client.Mutex.Unlock()
	}
	for _, deadClient := range deadClients {
		c.RemoveLocked(deadClient)
	}
	deadClients = deadClients[:0]
	MapMutex.Unlock()
}

// Publish sends a message to all clients that are subscribed to the given topic
func (c Clients) Publish(target uint8, data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		if client.Subscriptions.ToBitmask()&target != 0 {
			client.Mutex.Lock()
			err := client.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				deadClients = append(deadClients, client.SocketId)
			}
			client.Mutex.Unlock()
		}
	}
	for _, deadClient := range deadClients {
		c.RemoveLocked(deadClient)
	}
	deadClients = deadClients[:0]
	MapMutex.Unlock()
}

// PublishExcept sends a message to all clients that are subscribed to the given topic except the one with the given uuid
func (c Clients) PublishExcept(target uint8, socketId uint32, data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		if client.SocketId != socketId && client.Subscriptions.ToBitmask()&target != 0 {
			client.Mutex.Lock()
			err := client.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				deadClients = append(deadClients, client.SocketId)
			}
			client.Mutex.Unlock()
		}
	}
	for _, deadClient := range deadClients {
		c.RemoveLocked(deadClient)
	}
	deadClients = deadClients[:0]
	MapMutex.Unlock()
}

// BroadcastExcept sends a message to all clients except the one with the given uuid
func (c Clients) BroadcastExcept(socketId uint32, data []byte) {
	MapMutex.Lock()
	for _, client := range c {
		if client.SocketId != socketId {
			client.Mutex.Lock()
			err := client.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				deadClients = append(deadClients, client.SocketId)
			}
			client.Mutex.Unlock()
		}
	}
	for _, deadClient := range deadClients {
		c.RemoveLocked(deadClient)
	}
	deadClients = deadClients[:0]
	MapMutex.Unlock()
}

func (c Clients) Get(socketId uint32) *Client {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	client, ok := c[socketId]
	if !ok {
		return nil
	}
	return &client
}

func GetNumberOfClients() uint32 {
	return NbClients
}
