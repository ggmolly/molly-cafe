package watchdogs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorIdleUptime(packet *socket.Packet) {
	file, err := os.Open("/proc/uptime")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for {
		file.Seek(0, 0)
		var idleUptime, totalUptime float64
		_, err = fmt.Fscanf(file, "%f %f", &totalUptime, &idleUptime)
		if err != nil {
			log.Println("/!\\ Failed to parse uptime file", err)
		}
		packet.SetPercentage(float32(idleUptime) / float32(totalUptime))
		time.Sleep(REFRESH_DELAY)
	}
}
