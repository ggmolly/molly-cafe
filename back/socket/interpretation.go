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
	okayPackets[packets.MOUSE_MOVE_ID] = Requirements{Equal: 4, Function: nil} // expecting two float32
}

func validPacket(data []byte, reqs Requirements) bool {
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
		return false
	}
	// Check function requirements
	if reqs.Function != nil {
		return (*reqs.Function)(data)
	}
	return true
}

func SanitizePacket(data []byte) (uint8, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty packet?")
	}
	packetId := data[0]
	if _, ok := okayPackets[packetId]; !ok {
		return 0, fmt.Errorf("invalid packet type: %d", packetId)
	}
	if !validPacket(data[1:], okayPackets[packetId]) {
		return 0, fmt.Errorf("invalid packet size: %d", len(data))
	}
	return packetId, nil
}
