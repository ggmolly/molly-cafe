package routes

import (
	"github.com/bettercallmolly/illustrious/socket"
	"github.com/gofiber/contrib/websocket"
)

func WSRoutine(c *websocket.Conn) {
	var (
		mt  int    // message type
		_   []byte // message
		err error  // error
	)
	if c.Locals("allowed") == nil {
		c.Close()
		return
	}
	socketId := socket.GenerateClientId()
	socket.ConnectedClients.Add(c, socketId)
	defer func() { // Avoid resource leak
		socket.ConnectedClients.Remove(socketId)
		// Broadcast the disconnected client to all other clients
		socket.ConnectedClients.BroadcastExcept(socketId, []byte{0xFE})
	}()
	// Broadcast the new client to all other clients
	socket.ConnectedClients.BroadcastExcept(socketId, []byte{0xFF})
	count := socket.GetNumberOfClients() // uint32
	buffer := []byte{0xFD, byte(count >> 24), byte(count >> 16), byte(count >> 8), byte(count)}
	c.WriteMessage(websocket.BinaryMessage, buffer)
	// Send packets to the connected client
	for _, packet := range socket.PacketMap {
		c.WriteMessage(websocket.BinaryMessage, packet.GetRawBytes())
	}
	for {
		// Keep connection alive, and if a message is received, disconnect the client
		if mt, _, err = c.ReadMessage(); err != nil || mt == websocket.CloseMessage {
			return
		}
	}
}
