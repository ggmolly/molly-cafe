package routes

import (
	"bytes"
	"strings"

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
	Playing = false
)

const (
	InvalidReqBody = "Invalid request body"
	InvalidCover   = "Invalid cover image"
)

func publishNewSong(dto dtoUpdate) {
	packet := socket.PacketMap["strawberry"]
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

	packet.Data = dataBuffer.Bytes()
	// Send the packet to all connected clients
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func publishState(dto dtoState) {
	packet := socket.PacketMap["strawberryState"]
	packet.Data = []byte{0x00}
	if dto.Playing {
		packet.Data[0] = 0x01
	}
	// Send the packet to all connected clients
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func publishSeek(dto dtoSeek) {
	packet := socket.PacketMap["strawberrySeek"]
	// Send the packet to all connected clients
	var dataBuffer bytes.Buffer
	for i := 0; i < 4; i++ {
		dataBuffer.WriteByte(byte(dto.Position >> uint(24-i*8)))
	}
	packet.Data = dataBuffer.Bytes()
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
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
	return nil
}
