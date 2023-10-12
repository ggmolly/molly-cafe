package socket

import (
	"errors"
	"sync"
)

type T_PacketMap map[string]*Packet

var (
	PacketMap         = make(T_PacketMap)
	InvalidPacketErr  = errors.New("invalid packet size")
	ErrInvalidChannel = errors.New("invalid subscription channel")

	packetMapMutex = &sync.Mutex{}
)

// GetPacketByName returns a packet from the packet map by name
func (packetMap *T_PacketMap) GetPacketByName(name string) (*Packet, bool) {
	packetMapMutex.Lock()
	defer packetMapMutex.Unlock()
	packet, ok := PacketMap[name]
	return packet, ok
}

// GetPacketByNameUnsafe returns a packet from the packet map by name without returning a bool
func (packetMap *T_PacketMap) GetPacketByNameUnsafe(name string) *Packet {
	packetMapMutex.Lock()
	defer packetMapMutex.Unlock()
	packet, _ := PacketMap[name]
	return packet
}

// AddPacket adds a packet to the packet map
func (packetMap *T_PacketMap) AddPacket(name string, packet *Packet) {
	packetMapMutex.Lock()
	defer packetMapMutex.Unlock()
	PacketMap[name] = packet
}

// RemovePacket removes a packet from the packet map
func (packetMap *T_PacketMap) RemovePacket(name string) {
	packetMapMutex.Lock()
	defer packetMapMutex.Unlock()
	delete(PacketMap, name)
}
