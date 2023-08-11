package watchdogs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func getCPUPacket(packetMaps *map[string]*socket.Packet, label string) *socket.Packet {
	label = strings.ToLower(label)
	packet, ok := (*packetMaps)[label]
	if !ok {
		packet = socket.NewPacket(socket.C_MISC, socket.DT_TEMPERATURE, label)
		(*packetMaps)[label] = packet
		return packet
	}
	return packet
}

func MonitorCPUTemp(packetMaps *map[string]*socket.Packet, clients *socket.Clients) {
	var sensorDir string
	// Find all the sensor directories, and return the first k10temp or coretemp in /sys/class/hwmon
	dirEntry, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		panic(err)
	}
	for _, dir := range dirEntry {
		// try to open the name file
		data, readErr := os.ReadFile("/sys/class/hwmon/" + dir.Name() + "/name")
		if readErr != nil {
			break
		}
		if string(data) == "k10temp\n" || string(data) == "coretemp\n" {
			sensorDir = "/sys/class/hwmon/" + dir.Name()
			break
		}
	}

	// If we didn't find a sensor, just return
	if sensorDir == "" {
		log.Println("No CPU temperature sensor found")
		return
	}

	var paths []string
	// Scan sensor directory
	dirEntry, err = os.ReadDir(sensorDir)
	if err != nil {
		panic(err)
	}

	// Find all the temp*_input files
	for _, dir := range dirEntry {
		if strings.HasPrefix(dir.Name(), "temp") && strings.HasSuffix(dir.Name(), "input") {
			paths = append(paths, sensorDir+"/"+dir.Name())
		}
	}

	for {
		// Read all the temp*_input files and their corresponding temp*_label files
		for _, path := range paths {
			name, err := os.ReadFile(path[:len(path)-6] + "_label")
			if err != nil {
				panic(err)
			}
			val, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}

			// parse val into a float
			var temp float32
			_, err = fmt.Sscanf(string(val), "%f", &temp)
			if err != nil {
				panic(err)
			}
			packet := getCPUPacket(packetMaps, string(name[:len(name)-1])) // trim the newline
			packet.SetTemperature(temp / 1000)
		}
		time.Sleep(REFRESH_DELAY)
	}
}
