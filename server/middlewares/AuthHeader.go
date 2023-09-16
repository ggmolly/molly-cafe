package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func AuthHeader(c *fiber.Ctx) error {
	if c.Get("X-Mana-Key") != os.Getenv("API_KEY") {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	} else {
		return c.Next()
	}
}
