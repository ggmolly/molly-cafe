package socket

import (
	"bytes"
	"fmt"
	"math"
)

type Packet struct {
	Target     uint8
	Category   uint8
	Id         uint16
	DataType   uint8
	NameLength uint16
	Name       string
	Data       []byte
	Dirty      bool
}

const (
	// Monitoring stuff
	T_MONITORING = 0x00

	C_SERVICE       = 0x00
	C_HARD_RESOURCE = 0x01
	C_SOFT_RESOURCE = 0x02
	C_MISC          = 0x03

	DT_UINT8       = 0x00
	DT_UINT32      = 0x01
	DT_PERCENTAGE  = 0x02
	DT_TEMPERATURE = 0x03
	DT_LOAD_USAGE  = 0x04
	DT_SPECIAL     = 0xFF

	S_DEAD      = 0x00
	S_UNHEALTHY = 0x01
	S_HEALTHY   = 0x02
	S_OK        = 0x02

	// School projects stuff
	T_SCHOOL_PROJECTS = 0x01

	C_SCHOOL   = 0x00
	C_PERSONAL = 0x01

	P_WIP           = 0x00
	P_GRADING       = 0x01
	P_NOT_CONCERNED = 0xFF // mostly for personal projects

	// Pistache stuff
	T_PISTACHE = 0x02

	C_PISTACHE = 0x00

	// Strawberry stuff
	T_STRAWBERRY       = 0x03
	T_STRAWBERRY_SEEK  = 0x04
	T_STRAWBERRY_STATE = 0x05

	C_STRAWBERRY = 0x00

	// Cursor removal (0x06 -> MouseMovePacket)
	T_CURSOR_BYE = 0x07
	C_CURSOR_BYE = 0x00

	// Sleep tracking
	T_SLEEP = 0x08
	C_SLEEP = 0x00

	ERR_DATA_TYPE_MISMATCH = "dataType doesn't match the passed value"
)

var (
	packedId uint16 = 0
	// targetIdPrefix is used to generate the DOM element id for the packet
	targetIdPrefix = map[uint8]string{
		T_MONITORING:      "m",
		T_SCHOOL_PROJECTS: "sp",
		T_PISTACHE:        "pb",
	}
)

// Creates a new generic packet, this is the most common packet type
// and should be used in most cases (to send uint8, uint32, percentages...)
func NewPacket(target, category, dataType uint8, name string) *Packet {
	packedId++
	return &Packet{
		Target:   target,
		Category: category,
		Id:       packedId,
		DataType: dataType,
		Name:     name,
	}
}

// Creates a new packet, but the packet id stays at 0 and is not incremented
func NewUntrackedPacket(target, category, dataType uint8, name string) *Packet {
	return &Packet{
		Target:   target,
		Category: category,
		Id:       0,
		DataType: dataType,
		Name:     name,
	}
}

func NewMonitoringPacket(category, dataType uint8, name string) *Packet {
	return NewPacket(T_MONITORING, category, dataType, name)
}

// This function will remove the packet from the map, and remove the DOM element in the web page
func (p *Packet) RemoveDOM() {
	var buffer bytes.Buffer
	buffer.WriteByte(0xFC) // DOMPopPacket
	buffer.WriteString(fmt.Sprintf("%s-%d", targetIdPrefix[p.Target], p.Id))
	ConnectedClients.Broadcast(buffer.Bytes())
}

func (p *Packet) SetState(state uint8) *Packet {
	if p.DataType != DT_UINT8 {
		panic(ERR_DATA_TYPE_MISMATCH)
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
		panic(ERR_DATA_TYPE_MISMATCH)
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
		panic(ERR_DATA_TYPE_MISMATCH)
	}
	return p.setFloat32(value)
}

func (p *Packet) SetTemperature(value float32) *Packet {
	if p.DataType != DT_TEMPERATURE {
		panic(ERR_DATA_TYPE_MISMATCH)
	}
	return p.setFloat32(value)
}

func (p *Packet) SetLoadUsage(value float32) *Packet {
	if p.DataType != DT_LOAD_USAGE {
		panic(ERR_DATA_TYPE_MISMATCH)
	}
	return p.setFloat32(value)
}

func (p *Packet) GetRawBytes() []byte {
	buffer := make([]byte, 0, 8+len(p.Name)+len(p.Data))
	buffer = append(buffer, p.Target)
	buffer = append(buffer, p.Category)
	buffer = append(buffer, byte(p.Id>>8), byte(p.Id))
	buffer = append(buffer, p.DataType)
	buffer = append(buffer, byte(len(p.Name)>>8), byte(len(p.Name)))
	buffer = append(buffer, []byte(p.Name)...)
	buffer = append(buffer, p.Data...)
	return buffer
}
