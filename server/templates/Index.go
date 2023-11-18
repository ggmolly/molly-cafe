package templates

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/leitner"
	"github.com/bettercallmolly/illustrious/routes"
	"github.com/bettercallmolly/illustrious/socket"
	"github.com/bettercallmolly/illustrious/watchdogs"
	"github.com/gofiber/fiber/v2"
)

const (
	SleepGoal = 8 * 3600 // 8 hours
)

var (
	Birthday = time.Date(2003, 01, 02, 0, 0, 0, 0, time.UTC)
)

type projectPacket struct {
	watchdogs.SchoolProject
	ID uint16 `json:"id"`
}

type pistachePost struct {
	watchdogs.PistachePost
	ID uint16 `json:"id"`
}

func getProjects() []projectPacket {
	var projects []projectPacket
	// Loop through the packet map
	for _, packet := range socket.PacketMap {
		if packet.Target == socket.T_SCHOOL_PROJECTS {
			projects = append(projects, projectPacket{
				SchoolProject: watchdogs.DeserializeProject(packet),
				ID:            packet.Id,
			})
		}
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}

func getPistachePosts() []pistachePost {
	var posts []pistachePost
	// Loop through the packet map
	for _, packet := range socket.PacketMap {
		if packet.Target == socket.T_PISTACHE {
			posts = append(posts, pistachePost{
				PistachePost: watchdogs.DeserializePistachePost(packet),
				ID:           packet.Id,
			})
		}
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreationDate.After(posts[j].CreationDate)
	})
	return posts
}

func Index(c *fiber.Ctx) error {
	hiraganas, _ := leitner.LeitnerData.GetTopic("hiraganas")
	katakanas, _ := leitner.LeitnerData.GetTopic("katakanas")
	packet, ok := socket.PacketMap.GetPacketByName("sleepTracking")
	var timeSlept int
	if ok {
		// deserialize the packet
		for i := 0; i < 4; i++ {
			timeSlept |= int(packet.Data[i]) << (8 * (3 - i))
		}
	}
	weatherCondition := ""
	if len(watchdogs.CachedWeatherData.Data.Weather) > 0 {
		weatherCondition = strings.ToLower(watchdogs.CachedWeatherData.Data.Weather[0].Main)
	}
	return c.Render("index", fiber.Map{
		"age":                 uint8(time.Since(Birthday).Hours() / 24 / 365),
		"learnedHiraganas":    hiraganas.CompletedCards,
		"totalHiraganas":      hiraganas.Total,
		"percentageHiraganas": uint8(float32(hiraganas.CompletedCards) / float32(hiraganas.Total) * 100),
		"learnedKatakanas":    katakanas.CompletedCards,
		"totalKatakanas":      katakanas.Total,
		"percentageKatakanas": uint8(float32(katakanas.CompletedCards) / float32(katakanas.Total) * 100),
		"learningStreak":      leitner.LeitnerData.GetStreak(),
		"sunriseTime":         watchdogs.CachedWeatherData.Data.Sys.Sunrise,
		"sunsetTime":          watchdogs.CachedWeatherData.Data.Sys.Sunset,
		"currentTime":         time.Now().Unix(), // Server time is Europe/Paris
		"weatherCondition":    weatherCondition,
		"cloudiness":          watchdogs.CachedWeatherData.Data.Clouds.All,
		"humidity":            watchdogs.CachedWeatherData.Data.Main.Humidity,
		"feltTemperature":     fmt.Sprintf("%.2f", watchdogs.CachedWeatherData.Data.Main.FeelsLike),
		"windSpeed":           fmt.Sprintf("%.2f", watchdogs.CachedWeatherData.Data.Wind.Speed),
		"sleepTime":           timeSlept,
		"projects":            getProjects(),
		"pistachePosts":       getPistachePosts(),
		"strawberryCover":     routes.CurrentlyPlaying.Cover,
		"strawberryTitle":     routes.CurrentlyPlaying.Title,
		"strawberryArtists":   strings.Join(routes.CurrentlyPlaying.Artists, ", "),
		"strawberryLength":    routes.CurrentlyPlaying.Length,
		"connectedCount":      socket.NbClients,
	})
}
