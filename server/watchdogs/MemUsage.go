package watchdogs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorMemUsage(packet *socket.Packet) {
	file, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal("/!\\ Failed to open /proc/meminfo", err)
	}
	defer file.Close()
	for {
		file.Seek(0, 0)
		scanner := bufio.NewScanner(file)
		var memAvailable, memTotal float32
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "MemAvailable:") {
				line = strings.TrimSpace(line[13 : len(line)-2])
				_, err := fmt.Sscanf(line, "%f", &memAvailable)
				if err != nil {
					log.Println("/!\\ Failed to parse field MemAvailable /proc/meminfo", err)
					break
				}
			} else if strings.HasPrefix(line, "MemTotal:") {
				line = strings.TrimSpace(line[9 : len(line)-2])
				_, err := fmt.Sscanf(line, "%f", &memTotal)
				if err != nil {
					log.Println("/!\\ Failed to parse field MemTotal /proc/meminfo", err)
					break
				}
			}
			if memAvailable > 0 && memTotal > 0 {
				packet.SetLoadUsage(100 - (memAvailable / memTotal * 100))
				break
			}
		}
		time.Sleep(REFRESH_DELAY)
	}
}
