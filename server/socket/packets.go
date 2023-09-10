package socket

import "errors"

var (
	PacketMap         = make(map[string]*Packet)
	InvalidPacketErr  = errors.New("invalid packet size")
	ErrInvalidChannel = errors.New("invalid subscription channel")
)
