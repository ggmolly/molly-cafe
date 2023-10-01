package routes

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/gofiber/fiber/v2"
)

type LeitnerTopic struct {
	CompletedCards uint32 `json:"completed_cards"`
	Total          uint32 `json:"total"`
}

type LeitnerConfig struct {
	Topics    map[string]LeitnerTopic `json:"topics"`
	Streak    uint32                  `json:"streak"`
	UpdatedAt time.Time               `json:"updated_at"`
}

func (c *LeitnerConfig) Save() error {
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(config)
}

func (c *LeitnerConfig) Load() error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(c)
}

const (
	configPath = "leitner.json"
)

var (
	config LeitnerConfig
)

func updateStreakPacket(newStreak uint32) {
	packet, ok := socket.PacketMap["L_STREAK"]
	if !ok {
		packet = socket.NewPacket(
			socket.T_LEITNER,
			socket.C_LEITNER_STREAK,
			socket.DT_UINT32,
			"",
		)
		packet.Data = make([]byte, 4)
		socket.PacketMap["L_STREAK"] = packet
	}
	// Update the packet
	packet.Data[0] = byte(newStreak >> 24)
	packet.Data[1] = byte(newStreak >> 16)
	packet.Data[2] = byte(newStreak >> 8)
	packet.Data[3] = byte(newStreak)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

// Update / set packet, and broadcast to all clients
func updateLeitnerPacket(topic string) {
	innerName := "L_" + topic
	packet, ok := socket.PacketMap[innerName]
	if !ok {
		packet = socket.NewPacket(
			socket.T_LEITNER,
			socket.C_LEITNER,
			socket.DT_SPECIAL,
			innerName, // DOM element id
		)
		packet.Data = make([]byte, 8)
		socket.PacketMap[innerName] = packet
	}
	// Update the packet
	completed := config.Topics[topic].CompletedCards
	total := config.Topics[topic].Total
	// Pack completed and total into the buffer (uint32 -> 4 bytes)
	packet.Data[0] = byte(completed >> 24)
	packet.Data[1] = byte(completed >> 16)
	packet.Data[2] = byte(completed >> 8)
	packet.Data[3] = byte(completed)
	// Pack completed and total into the buffer (uint32 -> 4 bytes)
	packet.Data[4] = byte(total >> 24)
	packet.Data[5] = byte(total >> 16)
	packet.Data[6] = byte(total >> 8)
	packet.Data[7] = byte(total)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func UpdateLeitner(c *fiber.Ctx) error {
	topic, ok := config.Topics[c.Params("topic")]
	if !ok {
		return c.SendStatus(404)
	}
	if topic.CompletedCards >= topic.Total {
		return c.SendStatus(400)
	}
	config.Topics[c.Params("topic")] = LeitnerTopic{
		CompletedCards: topic.CompletedCards + 1,
		Total:          topic.Total,
	}
	// Save the config
	config.Save()
	// Update the packet
	updateLeitnerPacket(c.Params("topic"))
	return c.SendStatus(200)
}

func UpdateLeitnerStreak(c *fiber.Ctx) error {
	// Check how many days have passed since the last update
	daysPassed := int(time.Since(config.UpdatedAt).Hours() / 24)
	if daysPassed == 0 {
		return c.SendStatus(202)
	} else if daysPassed > 1 {
		// Reset the streak
		config.Streak = 1
	} else {
		config.Streak++
	}

	// If the streak was updated, update the streak packet too
	updateStreakPacket(config.Streak)
	config.UpdatedAt = time.Now()
	config.UpdatedAt = config.UpdatedAt.Add(-time.Duration(config.UpdatedAt.Hour()) * time.Hour)
	config.UpdatedAt = config.UpdatedAt.Add(-time.Duration(config.UpdatedAt.Minute()) * time.Minute)
	config.UpdatedAt = config.UpdatedAt.Add(-time.Duration(config.UpdatedAt.Second()) * time.Second)
	config.UpdatedAt = config.UpdatedAt.Add(-time.Duration(config.UpdatedAt.Nanosecond()) * time.Nanosecond)

	// Save the config
	config.Save()

	return c.SendStatus(200)
}

func init() {
	// check if the file leitner.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Leitner config does not exist. The Leitner API will not be available.")
	} else {
		// Load the config
		if err := config.Load(); err != nil {
			log.Println("Could not load leitner config. The Leitner API will not be available.")
		}
	}
	// Update the packets
	for topic := range config.Topics {
		updateLeitnerPacket(topic)
	}
	updateStreakPacket(config.Streak)
}
