package packets

// Returns a welcome packet with the given uuid
func WriteWelcomePacket(uuid string) []byte {
	data := append([]byte{WELCOME_PACKET_ID}, []byte(uuid)...)
	return data
}
