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
	"github.com/bettercallmolly/illustrious/watchdogs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/joho/godotenv"
)

type UpdateDetails struct {
	Name   string `json:"name"`
	Status uint8  `json:"status"` // 0 = dead, 1 = unhealthy, 2 = healthy
}

var (
	REFRESH_DELAY = 5 * time.Second
	ProjectPath   string
	PistacheRoot  = "./pistache"
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

	ProjectPath, err := configuration.GetRootPath("projects")
	if err != nil {
		log.Println("'projects' folder could not be found. project management will be disabled")
	} else {
		log.Printf("Monitoring projects from %s", ProjectPath)
		go watchdogs.MonitorSchoolProjects(&socket.PacketMap, filepath.Join(ProjectPath, "school"))
	}

	PistacheRoot, err = configuration.GetRootPath("pistache")
	if err != nil {
		log.Println("'pistache' folder could not be found. pistache will be disabled")
	} else {
		log.Printf("Pistache root set to %s", PistacheRoot)
		go watchdogs.MonitorPistachePosts(&socket.PacketMap, PistacheRoot)
	}
}

func main() {
	app := fiber.New(
		fiber.Config{
			IdleTimeout:             2 * time.Second,
			ProxyHeader:             "CF-Connecting-IP",
			EnableTrustedProxyCheck: true,
		},
	)

	app.Use(cache.New(cache.Config{
		Expiration:   1 * time.Hour,
		CacheControl: true,
	}))

	app.Use("/ws", middlewares.WebSocketUpgrade)

	app.Get("/ws", websocket.New(routes.WSRoutine))

	strawberryAPI := app.Group("/api/strawberry", middlewares.LANOnly)
	strawberryAPI.Post("/", routes.StrawberryUpdate)
	strawberryAPI.Patch("/seek", routes.SetStrawberrySeek)
	strawberryAPI.Patch("/state", routes.SetStrawberryState)

	// iOS shortcuts will send a POST request to this endpoint
	app.Post("/api/sleep", routes.SleepTracking)

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

	// Serve static files in argv[1]
	if (len(os.Args) < 2) || (os.Args[1] == "") {
		panic("No directory specified")
	}
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		panic("Directory does not exist")
	}
	app.Static("/", os.Args[1])
	app.Static("/pistache", PistacheRoot)
	app.Listen("0.0.0.0:50154")
}
