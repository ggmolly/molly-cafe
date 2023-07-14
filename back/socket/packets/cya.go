package packets

// Returns a cya packet with the given uuid
func WriteCyaPacket(socketId uint16) []byte {
	b1 := byte(socketId >> 8)
	b2 := byte(socketId)
	data := append([]byte{CYA_PACKET_ID}, b1, b2)
	return data
}
