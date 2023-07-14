package main

import (
	"log"
	"time"

	"github.com/bettercallmolly/molly/game"
	"github.com/bettercallmolly/molly/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
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
		socketId := socket.GenerateClientId()
		socket.ConnectedClients.Add(c, socketId)
		defer func() { // Avoid resource leak
			socket.ConnectedClients.Remove(socketId)
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
					game.HandlePacket(socketId, id, msg)
				} else { // Invalid message type, disconnect the client and break the loop
					log.Printf("Invalid message type: %d, disconnecting client", mt)
					socket.ConnectedClients.Remove(socketId)
					break
				}
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
