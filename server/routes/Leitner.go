package routes

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/gofiber/fiber/v2"
)

type LeitnerTopic struct {
	CompletedCards uint32 `json:"completed_cards"`
	Total          uint32 `json:"total"`
}

type LeitnerConfig struct {
	Topics map[string]LeitnerTopic `json:"topics"`
}

const (
	configPath = "leitner.json"
)

var (
	config LeitnerConfig
)

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
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return c.SendStatus(500)
	}
	defer file.Close()
	json.NewEncoder(file).Encode(config)
	// Update the packet
	updateLeitnerPacket(c.Params("topic"))
	return c.SendStatus(200)
}

func init() {
	// check if the file leitner.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Leitner config does not exist. The Leitner API will not be available.")
	} else {
		// if it exists, read the config
		file, _ := os.Open(configPath)
		defer file.Close()
		err := json.NewDecoder(file).Decode(&config)
		if err != nil {
			log.Println("Error decoding leitner config:", err)
		}
	}
}
