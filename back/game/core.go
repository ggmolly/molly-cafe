package game

import (
	"github.com/bettercallmolly/molly/game/actions"
	"github.com/bettercallmolly/molly/socket/packets"
)

var (
	dispatchers = make(map[uint8]func(string, []byte))
)

func init() {
	dispatchers[packets.MOUSE_MOVE_ID] = actions.HandleMouseMove
}

func HandlePacket(senderUUID string, packetId uint8, data []byte) {
	dispatchers[packetId](senderUUID, data)
}
