package leitner

import "github.com/bettercallmolly/illustrious/socket"

// Update / set packet, and broadcast to all clients
func UpdateLeitnerPacket(topic string) {
	innerName := "L_" + topic
	packet, ok := socket.PacketMap.GetPacketByName(innerName)
	if !ok {
		packet = socket.NewPacket(
			socket.T_LEITNER,
			socket.C_LEITNER,
			socket.DT_SPECIAL,
			innerName, // DOM element id
		)
		packet.Data = make([]byte, 8)
		socket.PacketMap.AddPacket(innerName, packet)
	}
	// Update the packet
	configMutex.Lock()
	completed := LeitnerData.Topics[topic].CompletedCards
	total := LeitnerData.Topics[topic].Total
	configMutex.Unlock()
	// Pack completed and total into the buffer (uint32 -> 4 bytes)
	packet.Data[0] = byte(completed >> 24)
	packet.Data[1] = byte(completed >> 16)
	packet.Data[2] = byte(completed >> 8)
	packet.Data[3] = byte(completed)
	// Pack completed and total into the buffer (uint32 -> 4 bytes)
	packet.Data[4] = byte(total >> 24)
	packet.Data[5] = byte(total >> 16)
	packet.Data[6] = byte(total >> 8)
	packet.Data[7] = byte(total)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}
