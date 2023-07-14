package socket

import (
	"fmt"

	"github.com/bettercallmolly/molly/socket/packets"
)

// A function that does nothing, called if a unexpected packet is received
func noop() error {
	return nil
}

var (
	PacketHandlers = map[uint8]func() error{}
)

func init() {
	PacketHandlers[packets.CYA_PACKET_ID] = noop
	PacketHandlers[packets.WELCOME_PACKET_ID] = noop
}

func InterpretPacket(data []byte) error {
	packetType := data[0]
	if handler, ok := PacketHandlers[packetType]; ok {
		return handler()
	} else {
		return fmt.Errorf("unknown packet type: %d", packetType)
	}
}
