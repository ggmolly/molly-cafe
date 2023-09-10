package packets

import (
	"fmt"
	"math"

	"github.com/bettercallmolly/illustrious/socket"
)

// CursorPositionPacket is a packet that contains the position of a cursor
func CursorPositionPacket(sender *socket.Client, payload []byte) error {
	if len(payload) != 9 {
		return socket.InvalidPacketErr
	}
	// Parse 4 bytes as float32 for x, 4 bytes as float32 for y
	x := math.Float32frombits(uint32(payload[1])<<24 | uint32(payload[2])<<16 | uint32(payload[3])<<8 | uint32(payload[4]))
	y := math.Float32frombits(uint32(payload[5])<<24 | uint32(payload[6])<<16 | uint32(payload[7])<<8 | uint32(payload[8]))
	socketId := sender.SocketId
	if x > 100 || y > 100 {
		return fmt.Errorf("invalid cursor position: %f, %f", x, y)
	}
	payload = append(payload, byte(socketId>>24), byte(socketId>>16), byte(socketId>>8), byte(socketId))
	socket.ConnectedClients.PublishExcept(socket.MASK_CURSORS, sender.SocketId, payload)
	return nil
}

// CursorByePacket is a packet that contains the id of a cursor that has disconnected / unsubscribed
// It is used to remove the cursor from the client's screen
func CursorByePacket(leaverId uint32) *socket.Packet {
	packet := socket.NewUntrackedPacket(
		socket.T_CURSOR_BYE,
		socket.C_CURSOR_BYE,
		socket.DT_UINT32,
		"",
	)
	packet.SetUint32(leaverId)
	return packet
}
