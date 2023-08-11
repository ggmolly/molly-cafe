package watchdogs

import (
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorRunningProcesses(packet *socket.Packet) {
	for {
		// Open directory
		dir, err := os.Open("/proc")
		if err != nil {
			log.Println("/!\\ Could not open /proc", err)
		}
		var n uint32
		// Read directory
		for {
			// Read directory entry
			entry, err := dir.Readdir(1)
			if err != nil {
				break
			}
			// Check if it's a directory
			if entry[0].IsDir() && entry[0].Name()[0] >= '0' && entry[0].Name()[0] <= '9' {
				n++
			}
		}
		// Close directory
		dir.Close()
		packet.SetUint32(n)
		time.Sleep(REFRESH_DELAY)
	}
}
