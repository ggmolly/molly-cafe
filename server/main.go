package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bettercallmolly/illustrious/configuration"
	"github.com/bettercallmolly/illustrious/middlewares"
	"github.com/bettercallmolly/illustrious/routes"
	"github.com/bettercallmolly/illustrious/socket"
	"github.com/bettercallmolly/illustrious/templates"
	"github.com/bettercallmolly/illustrious/watchdogs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

type UpdateDetails struct {
	Name   string `json:"name"`
	Status uint8  `json:"status"` // 0 = dead, 1 = unhealthy, 2 = healthy
}

var (
	REFRESH_DELAY = 5 * time.Second
)

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		// Panic because some endpoints won't be protected by auth if this fails
		panic("Failed to load .env file.")
	}
	if len(os.Getenv("API_KEY")) < 128 {
		panic("API_KEY must be at least 128 characters long")
	}
	// Load config file
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		log.Println("Failed to load the config.json file, using default values")
	}
	// Force to load paths
	configuration.LoadPaths()

	socket.ConnectedClients = socket.NewClients()

	// TCP / UDP connections
	socket.PacketMap["tcp"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "tcp connections")
	socket.PacketMap["udp"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "udp connections")
	go watchdogs.MonitorSocketConnections(socket.PacketMap["tcp"], socket.PacketMap["udp"])

	// Dirtymem
	socket.PacketMap["packetDirtyMem"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "dirty mem (kB)")
	go watchdogs.MonitorDirtyMem(socket.PacketMap["packetDirtyMem"])

	// Opened file descriptors
	socket.PacketMap["openFiles"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "opened fds")
	go watchdogs.MonitorOpenFiles(socket.PacketMap["openFiles"])

	// Running processes
	socket.PacketMap["proccesses"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "running proccesses")
	go watchdogs.MonitorRunningProcesses(socket.PacketMap["proccesses"])

	// Idle uptime
	socket.PacketMap["idleUptime"] = socket.NewMonitoringPacket(socket.C_SOFT_RESOURCE, socket.DT_PERCENTAGE, "idle uptime")
	go watchdogs.MonitorIdleUptime(socket.PacketMap["idleUptime"])

	// Users logged in (disabled until a more efficient way is used)
	// socket.PacketMap["usersLoggedIn"] = socket.NewMonitoringPacket(socket.C_MISC, socket.DT_UINT32, "users logged in")
	// go watchdogs.MonitorLoggedUsers(socket.PacketMap["usersLoggedIn"])

	// Containers / Services
	go watchdogs.MonitorContainers(&socket.PacketMap)
	go watchdogs.ManualServices(&socket.PacketMap)

	// CPU temperature
	go watchdogs.MonitorCPUTemp(&socket.PacketMap)

	// Internet Speed
	socket.PacketMap["downSpeed"] = socket.NewMonitoringPacket(socket.C_MISC, socket.DT_UINT32, "down speed (Mbps)")
	go watchdogs.MonitorInternetSpeed(socket.PacketMap["downSpeed"])

	// RAM usage
	socket.PacketMap["ramUsage"] = socket.NewMonitoringPacket(socket.C_HARD_RESOURCE, socket.DT_LOAD_USAGE, "ram usage")
	go watchdogs.MonitorMemUsage(socket.PacketMap["ramUsage"])

	// Disk usage
	go watchdogs.MonitorDiskSpace(&socket.PacketMap)

	// Weather
	go watchdogs.MonitorWeather(&socket.PacketMap)

	// Pistache
	go watchdogs.MonitorPistachePosts(&socket.PacketMap, configuration.PistacheRoot)

	// School projects
	go watchdogs.MonitorSchoolProjects(&socket.PacketMap, filepath.Join(configuration.ProjectPath, "school"))
}

func main() {
	engine := html.New(configuration.TemplateRoot, ".html")

	if os.Getenv("MODE") == "dev" {
		engine.Reload(true)
	}

	engine.AddFunc("formatTimestamp", templates.FormatTimestamp)
	engine.AddFunc("formatSeconds", templates.FormatSeconds)
	engine.AddFunc("formatDuration", templates.FormatDuration)
	engine.AddFunc("getSleepColor", templates.GetSleepColor)
	engine.AddFunc("getGradeColor", templates.GetGradeColor)
	engine.AddFunc("formatDate", templates.FormatDate)
	engine.AddFunc("timeToUnix", templates.TimeToUnix)

	app := fiber.New(
		fiber.Config{
			AppName:                 "molly's cafe",
			IdleTimeout:             2 * time.Second,
			ProxyHeader:             "CF-Connecting-IP",
			EnableTrustedProxyCheck: true,
			ReadTimeout:             time.Second * 10,
			WriteTimeout:            time.Second * 10,
			Views:                   engine,
		},
	)

	// app.Use(cache.New(cache.Config{
	// 	Expiration:   1 * time.Hour,
	// 	CacheControl: true,
	// }))

	app.Use("/ws", middlewares.WebSocketUpgrade)

	app.Get("/ws", websocket.New(routes.WSRoutine))

	strawberryAPI := app.Group("/api/strawberry", middlewares.LANOnly)
	strawberryAPI.Post("/", routes.StrawberryUpdate)
	strawberryAPI.Patch("/seek", routes.SetStrawberrySeek)
	strawberryAPI.Patch("/state", routes.SetStrawberryState)

	// iOS shortcuts will send a POST request to this endpoint
	app.Post("/api/sleep", routes.SleepTracking)

	// Leitner API
	leitnerAPI := app.Group("/api/leitner", middlewares.LANOnly)
	leitnerAPI.Post("/:topic", routes.UpdateLeitner)
	leitnerAPI.Patch("/streak", routes.UpdateLeitnerStreak)

	go func() {
		for {
			// Broadcast the number of clients to all clients
			for _, packet := range socket.PacketMap {
				if packet.Dirty {
					socket.ConnectedClients.Broadcast(packet.GetRawBytes())
					packet.Dirty = false
				}
			}
			time.Sleep(REFRESH_DELAY)
		}
	}()
	app.Get("/", templates.Index)
	app.Static("/assets", filepath.Join(configuration.TemplateRoot, "assets"))
	app.Static("/fonts", filepath.Join(configuration.TemplateRoot, "assets", "fonts"))
	app.Static("/pistache", configuration.PistacheRoot)
	app.Listen("0.0.0.0:50154")
}
