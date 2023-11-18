package routes

import (
	"log"
	"os"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/gofiber/fiber/v2"
)

var (
	TimeSlept uint32
)

func SleepTracking(c *fiber.Ctx) error {
	var data struct {
		Time uint32 `json:"time"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	log.Println("Received sleep tracking data:", data.Time, "seconds")
	packet, ok := socket.PacketMap.GetPacketByName("sleepTracking")
	if !ok {
		packet = socket.NewPacket(socket.T_SLEEP, socket.C_SLEEP, socket.DT_UINT32, "")
		socket.PacketMap.AddPacket("sleepTracking", packet)
	}
	packet.SetUint32(data.Time)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())

	// Save to a file to be able to display the last sleep tracking data if the server restarts
	err := os.WriteFile(".last_sleep", []byte{byte(data.Time >> 24), byte(data.Time >> 16), byte(data.Time >> 8), byte(data.Time)}, 0644)
	if err != nil {
		// Not very important, just log it
		log.Println("Failed to save sleep tracking data to file:", err)
	}
	return c.SendStatus(201)
}

func init() {
	// Load last sleep tracking data
	data, err := os.ReadFile(".last_sleep")
	if err != nil {
		return
	}
	packet, ok := socket.PacketMap.GetPacketByName("sleepTracking")
	if !ok {
		packet = socket.NewPacket(socket.T_SLEEP, socket.C_SLEEP, socket.DT_UINT32, "")
		socket.PacketMap.AddPacket("sleepTracking", packet)
	}
	TimeSlept = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	packet.SetUint32(TimeSlept)
}
