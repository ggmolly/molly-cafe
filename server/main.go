package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/middlewares"
	"github.com/bettercallmolly/illustrious/routes"
	"github.com/bettercallmolly/illustrious/socket"
	"github.com/bettercallmolly/illustrious/watchdogs"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

type UpdateDetails struct {
	Name   string `json:"name"`
	Status uint8  `json:"status"` // 0 = dead, 1 = unhealthy, 2 = healthy
}

type Configuration struct {
	MonitoredServices []string          `json:"services"`
	DiskTranslations  map[string]string `json:"disk_translations"`
}

var (
	REFRESH_DELAY = 5 * time.Second
	Config        = Configuration{}
)

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Println("Failed to read the config.json file")
		return
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("Failed to parse the config.json file")
		return
	}
	log.Println("Configuration loaded !")
	log.Println("Monitored services:")
	for _, service := range Config.MonitoredServices {
		log.Println("  -", service)
	}
	log.Println("Disk translations:")
	for key, value := range Config.DiskTranslations {
		log.Println("  -", key, "->", value)
	}
}

func pollConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to create a watcher for the config file")
		return
	}
	defer watcher.Close()
	err = watcher.Add("config.json")
	if err != nil {
		log.Println("Failed to add the config file to the watcher")
		return
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Failed to get the event from the watcher")
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Config file changed, reloading...")
				loadConfig()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("Failed to get the error from the watcher")
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func init() {
	// Load config file
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		log.Println("Failed to load the config.json file, using default values")
	} else {
		loadConfig()
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
	go watchdogs.ManualServices(&socket.PacketMap, Config.MonitoredServices...)

	// CPU temperature
	go watchdogs.MonitorCPUTemp(&socket.PacketMap)

	// Internet Speed
	socket.PacketMap["downSpeed"] = socket.NewMonitoringPacket(socket.C_MISC, socket.DT_UINT32, "down speed (Mbps)")
	go watchdogs.MonitorInternetSpeed(socket.PacketMap["downSpeed"])

	// RAM usage
	socket.PacketMap["ramUsage"] = socket.NewMonitoringPacket(socket.C_HARD_RESOURCE, socket.DT_LOAD_USAGE, "ram usage")
	go watchdogs.MonitorMemUsage(socket.PacketMap["ramUsage"])

	// Disk usage
	go watchdogs.MonitorDiskSpace(&socket.PacketMap, &Config.DiskTranslations)

	// Monitor config changes
	go pollConfig()
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
	app.Listen("127.0.0.1:50154")
}
