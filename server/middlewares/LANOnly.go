package middlewares

import (
	"net"

	"github.com/gofiber/fiber/v2"
)

func LANOnly(c *fiber.Ctx) error {
	// Drop if the request has CF-Connecting-IP header
	if c.Get("CF-Connecting-IP") != "" {
		return c.SendStatus(404)
	}
	// Parse the IP address
	ip := net.ParseIP(c.IP())
	// Drop if the IP address is not in the LAN range
	if !ip.IsLoopback() && !ip.IsPrivate() {
		return c.SendStatus(404)
	}
	return c.Next()
}
