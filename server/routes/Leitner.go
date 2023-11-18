package routes

import (
	"time"

	"github.com/bettercallmolly/illustrious/leitner"
	"github.com/gofiber/fiber/v2"
)

func UpdateLeitner(c *fiber.Ctx) error {
	topic, ok := leitner.LeitnerData.GetTopic(c.Params("topic"))
	if !ok {
		return c.SendStatus(404)
	}
	if topic.CompletedCards >= topic.Total {
		return c.SendStatus(400)
	}
	leitner.LeitnerData.SetTopic(c.Params("topic"), leitner.LeitnerTopic{
		CompletedCards: topic.CompletedCards + 1,
		Total:          topic.Total,
	})
	// Update the packet
	leitner.UpdateLeitnerPacket(c.Params("topic"))
	return c.SendStatus(200)
}

func UpdateLeitnerStreak(c *fiber.Ctx) error {
	// Check how many days have passed since the last update
	daysPassed := int(time.Since(leitner.LeitnerData.GetUpdatedAt()).Hours() / 24)
	if daysPassed == 0 {
		return c.SendStatus(202)
	} else if daysPassed > 1 {
		// Reset the streak
		leitner.LeitnerData.Streak = 1
	} else {
		leitner.LeitnerData.Streak++
	}

	// If the streak was updated, update the streak packet too
	leitner.UpdateStreakPacket(leitner.LeitnerData.Streak)
	leitner.LeitnerData.RefreshUpdatedAt()

	// Save the config
	leitner.LeitnerData.Save()

	return c.SendStatus(200)
}
