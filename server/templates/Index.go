package templates

import (
	"time"

	"github.com/bettercallmolly/illustrious/leitner"
	"github.com/gofiber/fiber/v2"
)

var (
	Birthday = time.Date(2003, 01, 02, 0, 0, 0, 0, time.UTC)
)

func Index(c *fiber.Ctx) error {
	hiraganas, _ := leitner.LeitnerData.GetTopic("hiraganas")
	katakanas, _ := leitner.LeitnerData.GetTopic("katakanas")
	return c.Render("index", fiber.Map{
		"age":                 uint8(time.Since(Birthday).Hours() / 24 / 365),
		"learnedHiraganas":    hiraganas.CompletedCards,
		"totalHiraganas":      hiraganas.Total,
		"percentageHiraganas": uint8(float32(hiraganas.CompletedCards) / float32(hiraganas.Total) * 100),
		"learnedKatakanas":    katakanas.CompletedCards,
		"totalKatakanas":      katakanas.Total,
		"percentageKatakanas": uint8(float32(katakanas.CompletedCards) / float32(katakanas.Total) * 100),
		"learningStreak":      leitner.LeitnerData.GetStreak(),
	})
}
