package socket

import (
	"bytes"
	"math"
)

type Packet struct {
	Category   uint8
	Id         uint16
	DataType   uint8
	NameLength uint16
	Name       string
	Data       []byte
	Dirty      bool
}

const (
	C_SERVICE       = 0x00
	C_HARD_RESOURCE = 0x01
	C_SOFT_RESOURCE = 0x02
	C_MISC          = 0x03

	DT_UINT8       = 0x00
	DT_UINT32      = 0x01
	DT_PERCENTAGE  = 0x02
	DT_TEMPERATURE = 0x03
	DT_LOAD_USAGE  = 0x04

	S_DEAD      = 0x00
	S_UNHEALTHY = 0x01
	S_HEALTHY   = 0x02
	S_OK        = 0x02
)

var (
	packedId uint16 = 0
)

func NewPacket(category uint8, dataType uint8, name string) *Packet {
	packedId++
	return &Packet{
		Category: category,
		Id:       packedId,
		DataType: dataType,
		Name:     name,
	}
}

func (p *Packet) SetState(state uint8) *Packet {
	if p.DataType != DT_UINT8 {
		panic("dataType does not match the passed value")
	}
	newData := []byte{state}
	if !bytes.Equal(newData, p.Data) {
		p.Data = newData
		p.Dirty = true
	}
	return p
}

func (p *Packet) setFloat32(value float32) *Packet {
	bits := math.Float32bits(value)
	newData := []byte{byte(bits >> 24), byte(bits >> 16), byte(bits >> 8), byte(bits)}
	if !bytes.Equal(newData, p.Data) {
		p.Data = newData
		p.Dirty = true
	}
	return p
}

func (p *Packet) SetUint32(value uint32) *Packet {
	if p.DataType != DT_UINT32 {
		panic("dataType does not match the passed value")
	}
	newData := []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)}
	if !bytes.Equal(newData, p.Data) {
		p.Data = newData
		p.Dirty = true
	}
	return p
}

func (p *Packet) SetPercentage(value float32) *Packet {
	if p.DataType != DT_PERCENTAGE {
		panic("dataType does not match the passed value")
	}
	return p.setFloat32(value)
}

func (p *Packet) SetTemperature(value float32) *Packet {
	if p.DataType != DT_TEMPERATURE {
		panic("dataType does not match the passed value")
	}
	return p.setFloat32(value)
}

func (p *Packet) SetLoadUsage(value float32) *Packet {
	if p.DataType != DT_LOAD_USAGE {
		panic("dataType does not match the passed value")
	}
	return p.setFloat32(value)
}

func (p *Packet) GetRawBytes() []byte {
	buffer := make([]byte, 0, 7+len(p.Name)+len(p.Data))
	buffer = append(buffer, p.Category)
	buffer = append(buffer, byte(p.Id>>8), byte(p.Id))
	buffer = append(buffer, p.DataType)
	buffer = append(buffer, byte(len(p.Name)>>8), byte(len(p.Name)))
	buffer = append(buffer, []byte(p.Name)...)
	buffer = append(buffer, p.Data...)
	return buffer
}
