package actions

import (
	"github.com/bettercallmolly/molly/socket"
	"github.com/bettercallmolly/molly/socket/packets"
)

func HandleMouseMove(socketId uint16, data []byte) {
	// Create a byte array of size 2 + 1 + 4
	// 1 byte for the packet ID
	// 2 bytes for the Socket ID
	// 2 bytes for the X coordinate
	// 2 bytes for the Y coordinate
	const size = 1 + 2 + 4
	payload := make([]byte, size)
	payload[0] = packets.MOUSE_MOVE_ID
	b1 := byte(socketId >> 8)
	b2 := byte(socketId)
	payload[1] = b1
	payload[2] = b2
	copy(payload[3:], data[1:])
	socket.ConnectedClients.BroadcastExcept(socketId, payload)
}
