package socket

import (
	"fmt"

	"github.com/bettercallmolly/molly/socket/packets"
)

type Requirements struct {
	Min      int
	Max      int
	Equal    int // If this is set, Min and Max are ignored
	Function *func([]byte) bool
}

var (
	okayPackets = map[uint8]Requirements{}
)

func init() {
	okayPackets[packets.MOUSE_MOVE_ID] = Requirements{Equal: 9, Function: nil} // expecting two float32
}

func validPacket(data []byte, reqs Requirements) (bool, bool) {
	// Check size requirements
	length := len(data)
	var lengthOk bool = true
	if reqs.Equal != 0 {
		lengthOk = length == reqs.Equal
	} else {
		if reqs.Min != 0 {
			lengthOk = length >= reqs.Min
		}
		if reqs.Max != 0 {
			lengthOk = length <= reqs.Max
		}
	}
	if !lengthOk {
		return false, false
	}
	// Check function requirements
	if reqs.Function != nil {
		return (*reqs.Function)(data), true
	}
	return true, true
}

func SanitizePacket(data []byte) (uint8, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty packet?")
	}
	packetId := data[0]
	if _, ok := okayPackets[packetId]; !ok {
		return 0, fmt.Errorf("invalid packet type: %d", packetId)
	}
	lengthOk, functionOk := validPacket(data, okayPackets[packetId])
	if !lengthOk {
		return 0, fmt.Errorf("invalid packet length: %d, expected %d", len(data), okayPackets[packetId].Equal)
	}
	if !functionOk {
		return 0, fmt.Errorf("invalid packet data: %v", data)
	}
	return packetId, nil
}
