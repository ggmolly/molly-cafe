package routes

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/gofiber/fiber/v2"
)

type dtoUpdate struct {
	Title   string   `json:"title"`
	Artists []string `json:"artists"`
	Cover   string   `json:"cover"`
	Length  uint32   `json:"length"`
}

type dtoSeek struct {
	Position uint32 `json:"position"`
}

type dtoState struct {
	Playing bool `json:"playing"`
}

var (
	Playing            = false
	CurrentTime uint32 = 0
	TimeMutex          = &sync.Mutex{}

	SeekPacket socket.Packet
)

const (
	InvalidReqBody = "Invalid request body"
	InvalidCover   = "Invalid cover image"
)

func getStrawberryPacket() *socket.Packet {
	packet, ok := socket.PacketMap.GetPacketByName("strawberry")
	if !ok {
		packet = socket.NewPacket(
			socket.T_STRAWBERRY,
			socket.C_STRAWBERRY,
			socket.DT_SPECIAL,
			"",
		)
		socket.PacketMap.AddPacket("strawberry", packet)
	}
	return packet
}

func updateTime(newTime uint32) {
	packet := getStrawberryPacket()
	length := len(packet.Data)
	if length < 4 {
		packet.Data = make([]byte, 4)
		length = 4
	}
	packet.Data[length-4] = byte(CurrentTime >> 24)
	packet.Data[length-3] = byte(CurrentTime >> 16)
	packet.Data[length-2] = byte(CurrentTime >> 8)
	packet.Data[length-1] = byte(CurrentTime)
}

func publishNewSong(dto dtoUpdate) {
	// Update the current time
	TimeMutex.Lock()
	Playing = true
	CurrentTime = 0
	TimeMutex.Unlock()

	packet := getStrawberryPacket()
	var dataBuffer bytes.Buffer

	separatedArtists := strings.Join(dto.Artists, ",")
	var titleLength uint16 = uint16(len(dto.Title))
	dataBuffer.WriteByte(byte(titleLength >> 8))
	dataBuffer.WriteByte(byte(titleLength))
	dataBuffer.WriteString(dto.Title)

	var artistsLength uint16 = uint16(len(separatedArtists))
	dataBuffer.WriteByte(byte(artistsLength >> 8))
	dataBuffer.WriteByte(byte(artistsLength))
	dataBuffer.WriteString(separatedArtists)

	var coverLength uint32 = uint32(len(dto.Cover))
	dataBuffer.WriteByte(byte(coverLength >> 24))
	dataBuffer.WriteByte(byte(coverLength >> 16))
	dataBuffer.WriteByte(byte(coverLength >> 8))
	dataBuffer.WriteByte(byte(coverLength))
	dataBuffer.WriteString(dto.Cover)

	// Write song length
	for i := 0; i < 4; i++ {
		dataBuffer.WriteByte(byte(dto.Length >> uint(24-i*8)))
	}

	// Write current time
	for i := 0; i < 4; i++ {
		dataBuffer.WriteByte(byte(CurrentTime >> uint(24-i*8)))
	}

	packet.Data = dataBuffer.Bytes()
	// Send the packet to all connected clients
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func publishState(dto dtoState) {
	packet, ok := socket.PacketMap.GetPacketByName("strawberryState")
	if !ok {
		packet = socket.NewPacket(
			socket.T_STRAWBERRY_STATE,
			socket.C_STRAWBERRY,
			socket.DT_SPECIAL,
			"",
		)
		socket.PacketMap.AddPacket("strawberryState", packet)
	}
	packet.Data = []byte{0x00}
	if dto.Playing {
		packet.Data[0] = 0x01
	}
	// Send the packet to all connected clients
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func publishSeek(dto dtoSeek) {
	// Send the packet to all connected clients
	var dataBuffer bytes.Buffer
	for i := 0; i < 4; i++ {
		dataBuffer.WriteByte(byte(dto.Position >> uint(24-i*8)))
	}
	updateTime(CurrentTime)
	SeekPacket.Data = dataBuffer.Bytes()
	socket.ConnectedClients.Broadcast(SeekPacket.GetRawBytes())
}

func StrawberryUpdate(c *fiber.Ctx) error {
	var dto dtoUpdate
	if err := c.BodyParser(&dto); err != nil {
		c.Status(400).SendString(InvalidReqBody)
		return err
	}
	// Publish the new song
	go publishNewSong(dto)
	return nil
}

func SetStrawberrySeek(c *fiber.Ctx) error {
	var dto dtoSeek
	if err := c.BodyParser(&dto); err != nil {
		c.Status(400).SendString(InvalidReqBody)
		return err
	}
	TimeMutex.Lock()
	CurrentTime = dto.Position
	TimeMutex.Unlock()
	publishSeek(dto)
	return nil
}

func SetStrawberryState(c *fiber.Ctx) error {
	var dto dtoState
	if err := c.BodyParser(&dto); err != nil {
		c.Status(400).SendString(InvalidReqBody)
		return err
	}
	publishState(dto)
	TimeMutex.Lock()
	Playing = dto.Playing
	TimeMutex.Unlock()
	return nil
}

func init() {
	go func() {
		// Increment the current time every second, if the song is playing
		for {
			time.Sleep(1 * time.Second)
			if !Playing {
				continue
			}
			TimeMutex.Lock()
			CurrentTime += 1000000
			updateTime(CurrentTime)
			TimeMutex.Unlock()
		}
	}()
	SeekPacket = *socket.NewPacket(socket.T_STRAWBERRY_SEEK, socket.C_STRAWBERRY, socket.DT_SPECIAL, "")
}
