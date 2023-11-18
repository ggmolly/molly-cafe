package leitner

import "github.com/bettercallmolly/illustrious/socket"

func UpdateStreakPacket(newStreak uint32) {
	packet, ok := socket.PacketMap.GetPacketByName("L_STREAK")
	if !ok {
		packet = socket.NewPacket(
			socket.T_LEITNER,
			socket.C_LEITNER_STREAK,
			socket.DT_UINT32,
			"",
		)
		packet.Data = make([]byte, 4)
		socket.PacketMap.AddPacket("L_STREAK", packet)
	}
	// Update the packet
	packet.Data[0] = byte(newStreak >> 24)
	packet.Data[1] = byte(newStreak >> 16)
	packet.Data[2] = byte(newStreak >> 8)
	packet.Data[3] = byte(newStreak)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}
