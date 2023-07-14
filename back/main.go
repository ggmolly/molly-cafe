package main

import (
	"log"

	"github.com/bettercallmolly/molly/socket"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

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
					err := socket.InterpretPacket(msg)
					if err != nil {
						log.Printf("Error interpreting packet: %s", err)
					}
				} else { // Invalid message type, disconnect the client and break the loop
					socket.ConnectedClients.Remove(uuid)
					break
				}
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
