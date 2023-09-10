package routes

import (
	"fmt"
	"sync"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/bettercallmolly/illustrious/socket/packets"
	"github.com/gofiber/contrib/websocket"
)

const (
	// PT = Packet Type
	PT_SUBSCRIBE   = 0x00
	PT_UNSUBSCRIBE = 0x01
	PT_CURSOR      = 0x06
)

var (
	packetFnMap = map[byte]func(*socket.Client, []byte) error{
		PT_CURSOR: packets.CursorPositionPacket,
	}
	ReceivedPackets     = 0
	ReceivedPacketsSize = 0
	MutexTemp           = &sync.Mutex{}
)

func parsePacket(client socket.Client, payload []byte, mt int) error {
	if len(payload) < 2 || mt != websocket.BinaryMessage {
		return fmt.Errorf("invalid packet received")
	}
	switch payload[0] {
	case PT_SUBSCRIBE:
		return client.Subscribe(payload[1])
	case PT_UNSUBSCRIBE:
		if payload[1] == socket.S_CURSORS {
			// Send a bye packet to all clients except the one that left
			socket.ConnectedClients.PublishExcept(
				socket.MASK_CURSORS,
				client.SocketId,
				packets.CursorByePacket(client.SocketId).GetRawBytes(),
			)
		}
		return client.Unsubscribe(payload[1])
	default:
		if fn, ok := packetFnMap[payload[0]]; ok {
			return fn(&client, payload)
		}
		return fmt.Errorf("invalid packet received")
	}
}

func WSRoutine(c *websocket.Conn) {
	if c.Locals("allowed") == nil {
		c.Close()
		return
	}
	socketId := socket.GenerateClientId()
	socket.ConnectedClients.Add(c, socketId)
	defer func() { // Avoid resource leak
		// Broadcast the disconnected client to all other clients
		socket.ConnectedClients.BroadcastExcept(socketId, []byte{0xFE})
		// Check if the client was subscribed to cursors, if so, publish a cursor deletion packet
		if socket.ConnectedClients[socketId].Subscriptions.Cursors {
			socket.ConnectedClients.PublishExcept(
				socket.MASK_CURSORS,
				socketId,
				packets.CursorByePacket(socketId).GetRawBytes(),
			)
		}
		socket.ConnectedClients.Remove(socketId)
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
		mt, p, err := c.ReadMessage()
		if err != nil {
			break
		}
		if err := parsePacket(socket.ConnectedClients[socketId], p, mt); err != nil {
			// log.Println(err)
			break
		}
	}
}
