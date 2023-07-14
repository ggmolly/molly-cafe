package actions

import (
	"github.com/bettercallmolly/molly/socket"
	"github.com/bettercallmolly/molly/socket/packets"
)

func HandleMouseMove(senderUUID string, data []byte) {
	// Create a byte array of size 36 + 1 + 8
	// 36 bytes for the UUID
	// 1 byte for the packet ID
	// 4 bytes for the X coordinate
	// 4 bytes for the Y coordinate
	payload := make([]byte, 36+1+8)
	payload[0] = packets.MOUSE_MOVE_ID
	copy(payload[1:37], []byte(senderUUID))
	copy(payload[37:45], data)
	socket.ConnectedClients.Broadcast(senderUUID, payload)
}
