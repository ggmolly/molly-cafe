package socket

import (
	"fmt"

	"github.com/bettercallmolly/molly/socket/packets"
)

// A function that does nothing, called if a unexpected packet is received
func noop() (*GenericPacket, error) {
	return nil, nil
}

var (
	PacketHandlers = map[uint8]func() (*GenericPacket, error){}
)

type GenericPacket struct {
	Type uint8
	Data []byte
}

func init() {
	// We're not supposed to receive these packets from the client
	PacketHandlers[packets.CYA_PACKET_ID] = noop
	PacketHandlers[packets.WELCOME_PACKET_ID] = noop
}

func InterpretPacket(data []byte) (*GenericPacket, error) {
	packetType := data[0]
	if handler, ok := PacketHandlers[packetType]; ok {
		return handler()
	} else {
		return nil, fmt.Errorf("unknown packet type: %d", packetType)
	}
}
