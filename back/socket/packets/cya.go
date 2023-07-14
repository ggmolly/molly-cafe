package packets

// Returns a cya packet with the given uuid
func WriteCyaPacket(uuid string) []byte {
	data := append([]byte{CYA_PACKET_ID}, []byte(uuid)...)
	return data
}
