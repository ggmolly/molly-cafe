package watchdogs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorIdleUptime(packet *socket.Packet) {
	for {
		file, err := os.Open("/proc/uptime")
		if err != nil {
			log.Println("/!\\ Failed to open /proc/uptime")
			time.Sleep(REFRESH_DELAY)
			continue
		}
		var idleUptime, totalUptime float64
		_, err = fmt.Fscanf(file, "%f %f", &totalUptime, &idleUptime)
		if err != nil {
			log.Println("/!\\ Failed to open /proc/uptime")
			time.Sleep(REFRESH_DELAY)
			continue
		}
		packet.SetPercentage(float32(idleUptime) / float32(totalUptime))
		file.Close()
		time.Sleep(REFRESH_DELAY)
	}
}
