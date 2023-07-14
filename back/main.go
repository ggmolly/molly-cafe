package main

import (
	"log"
	"time"

	"github.com/bettercallmolly/molly/game"
	"github.com/bettercallmolly/molly/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	receivedPackets = 0
	bandwidth       = 0
)

func init() {
	socket.ConnectedClients = socket.NewClients()
}

func main() {
	app := fiber.New()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			log.Printf("Received %d packets in the last second (%.3f KB/s)", receivedPackets, float64(bandwidth)/1024)
			receivedPackets = 0
			bandwidth = 0
		}
	}()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true) // upgrade
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int    // message type
			msg []byte // message
			err error  // error
		)
		if c.Locals("allowed") == nil {
			c.Close()
			return
		}
		uuid := uuid.New().String()
		socket.ConnectedClients.Add(c, uuid)
		defer func() { // Avoid resource leak
			socket.ConnectedClients.Remove(uuid)
		}()
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			} else {
				// binary message only
				if mt == websocket.BinaryMessage {
					id, err := socket.SanitizePacket(msg)
					if err != nil {
						log.Printf("Error interpreting packet: %s", err)
						continue
					}
					receivedPackets++
					bandwidth += len(msg)
					game.HandlePacket(uuid, id, msg)
				} else { // Invalid message type, disconnect the client and break the loop
					log.Printf("Invalid message type: %d, disconnecting client", mt)
					socket.ConnectedClients.Remove(uuid)
					break
				}
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
