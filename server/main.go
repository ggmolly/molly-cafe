package main

import (
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/bettercallmolly/illustrious/watchdogs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

type UpdateDetails struct {
	Name   string `json:"name"`
	Status uint8  `json:"status"` // 0 = dead, 1 = unhealthy, 2 = healthy
}

var (
	packets       = make(map[string]*socket.Packet)
	REFRESH_DELAY = 5 * time.Second
)

func init() {
	socket.ConnectedClients = socket.NewClients()

	// TCP / UDP connections
	packets["tcp"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "tcp connections")
	packets["udp"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "udp connections")
	go watchdogs.MonitorSocketConnections(packets["tcp"], packets["udp"])

	// Dirtymem
	packets["packetDirtyMem"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "dirty mem (kB)")
	go watchdogs.MonitorDirtyMem(packets["packetDirtyMem"])

	// Opened file descriptors
	packets["openFiles"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "opened fds")
	go watchdogs.MonitorOpenFiles(packets["openFiles"])

	// Running processes
	packets["proccesses"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_UINT32, "running proccesses")
	go watchdogs.MonitorRunningProcesses(packets["proccesses"])

	// Idle uptime
	packets["idleUptime"] = socket.NewPacket(socket.C_SOFT_RESOURCE, socket.DT_PERCENTAGE, "idle uptime")
	go watchdogs.MonitorIdleUptime(packets["idleUptime"])

	// Users logged in (disabled until a more efficient way is used)
	// packets["usersLoggedIn"] = socket.NewPacket(socket.C_MISC, socket.DT_UINT32, "users logged in")
	// go watchdogs.MonitorLoggedUsers(packets["usersLoggedIn"])

	// Containers / Services
	go watchdogs.MonitorContainers(&packets)
	go watchdogs.ManualServices(&packets, "nginx", "mariadb", "docker", "cron", "smbd")

	// CPU temperature
	go watchdogs.MonitorCPUTemp(&packets)

	// Internet Speed
	packets["downSpeed"] = socket.NewPacket(socket.C_MISC, socket.DT_UINT32, "down speed (Mbps)")
	go watchdogs.MonitorInternetSpeed(packets["downSpeed"])

	// RAM usage
	packets["ramUsage"] = socket.NewPacket(socket.C_HARD_RESOURCE, socket.DT_LOAD_USAGE, "ram usage")
	go watchdogs.MonitorMemUsage(packets["ramUsage"])

	// Disk usage
	go watchdogs.MonitorDiskSpace(&packets)
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
		Expiration:   8 * time.Hour,
		CacheControl: true,
	}))

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
			_   []byte // message
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
			// Broadcast the disconnected client to all other clients
			socket.ConnectedClients.BroadcastExcept(socketId, []byte{0xFE})
		}()
		// Broadcast the new client to all other clients
		socket.ConnectedClients.BroadcastExcept(socketId, []byte{0xFF})
		count := socket.GetNumberOfClients() // uint32
		buffer := []byte{0xFD, byte(count >> 24), byte(count >> 16), byte(count >> 8), byte(count)}
		c.WriteMessage(websocket.BinaryMessage, buffer)
		// Send packets to the connected client
		for _, packet := range packets {
			c.WriteMessage(websocket.BinaryMessage, packet.GetRawBytes())
		}
		for {
			// Keep connection alive, and if a message is received, disconnect the client
			if mt, _, err = c.ReadMessage(); err != nil || mt == websocket.CloseMessage {
				return
			}
		}
	}))

	go func() {
		for {
			// Broadcast the number of clients to all clients
			for _, packet := range packets {
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
