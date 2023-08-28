package watchdogs

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

var (
	diskTranslations = map[string]string{}
	translationsSet  = false
)

func getTranslatedName(path string) string {
	baseName := filepath.Base(path)
	if val, ok := diskTranslations[baseName]; ok {
		return val
	}
	return baseName
}

func getDiskPacket(packetMaps *map[string]*socket.Packet, path string) *socket.Packet {
	path = strings.ToLower(getTranslatedName(path))
	packet, ok := (*packetMaps)[path]
	if !ok {
		packet = socket.NewMonitoringPacket(socket.C_HARD_RESOURCE, socket.DT_LOAD_USAGE, path)
		(*packetMaps)[path] = packet
		return packet
	}
	return packet
}

// TODO: Optimize this
func MonitorDiskSpace(packetMaps *map[string]*socket.Packet, translations *map[string]string) {
	if !translationsSet {
		diskTranslations = *translations
		translationsSet = true
	}
	mountPoints := []string{}
	mounts, err := os.OpenFile("/proc/mounts", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(mounts)
	for scanner.Scan() {
		// line must start with /dev/sd
		if scanner.Text()[0:7] == "/dev/sd" && !strings.Contains(scanner.Text(), "docker") {
			slices := strings.Split(scanner.Text(), " ")
			mountPoints = append(mountPoints, slices[1])
			slices = nil // clear slices, we want this to be garbage collected asap
		}
	}
	mounts.Close()
	var stat syscall.Statfs_t
	for {
		for _, mountPoint := range mountPoints {
			err := syscall.Statfs(mountPoint, &stat)
			if err != nil {
				log.Println("/!\\ Unable to stat", mountPoint, err)
				continue
			}
			available := stat.Bavail * uint64(stat.Bsize)
			total := stat.Blocks * uint64(stat.Bsize)
			used := 100 - (float64(available) / float64(total) * 100)
			packet := getDiskPacket(packetMaps, mountPoint)
			packet.SetLoadUsage(float32(used))
		}
		time.Sleep(REFRESH_DELAY)
	}
}
