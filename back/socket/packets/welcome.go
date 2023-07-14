package packets

// Returns a welcome packet with the given uuid
func WriteWelcomePacket(socketId uint16) []byte {
	b1 := byte(socketId >> 8)
	b2 := byte(socketId)
	data := append([]byte{WELCOME_PACKET_ID}, b1, b2)
	return data
}
