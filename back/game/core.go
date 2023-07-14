package game

import (
	"github.com/bettercallmolly/molly/game/actions"
	"github.com/bettercallmolly/molly/socket/packets"
)

var (
	dispatchers = make(map[uint8]func(uint16, []byte))
)

func init() {
	dispatchers[packets.MOUSE_MOVE_ID] = actions.HandleMouseMove
}

func HandlePacket(socketId uint16, packetId uint8, data []byte) {
	dispatchers[packetId](socketId, data)
}
