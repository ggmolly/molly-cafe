package watchdogs

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

const (
	MEM_MULTIPLIER = 1024 // kB
)

func parseKb(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.ParseUint(s, 10, 32)
	return n * MEM_MULTIPLIER, err
}

func MonitorDirtyMem(dirtyMemPacket *socket.Packet) {
	for {
		file, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0)
		if err != nil {
			log.Println("/!\\ Failed to open /proc/meminfo", err)
			time.Sleep(REFRESH_DELAY)
			continue
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Dirty:") {
				dirtyMem, err := parseKb(strings.TrimSpace(line[6 : len(line)-3]))
				if err != nil {
					log.Println("/!\\ Failed to parse dirty memory", err)
				} else {
					dirtyMemPacket.SetUint32(uint32(dirtyMem / 1024))
				}
				break
			}
		}
		file.Close()
		time.Sleep(REFRESH_DELAY)
	}
}
