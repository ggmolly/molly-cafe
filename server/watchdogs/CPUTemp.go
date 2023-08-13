package watchdogs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

type Sensor struct {
	Name  string
	Value float32
	// This file won't be closed until the program exits, which is sad, but
	// a watchdog panics or the service gets stopped so this is fine I guess
	File *os.File
}

func (s *Sensor) Read() float32 {
	_, err := s.File.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fscanf(s.File, "%f", &s.Value)
	if err != nil {
		panic(err)
	}
	return s.Value
}

func getCPUPacket(packetMaps *map[string]*socket.Packet, label string) *socket.Packet {
	label = strings.ToLower(label)
	packet, ok := (*packetMaps)[label]
	if !ok {
		packet = socket.NewMonitoringPacket(socket.C_MISC, socket.DT_TEMPERATURE, label)
		(*packetMaps)[label] = packet
		return packet
	}
	return packet
}

func getSensorDir() (string, bool) {
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
			return "/sys/class/hwmon/" + dir.Name(), true
		}
	}
	return "", false
}

func getSensor(sensorDir string, name string, sensors *[]Sensor) {
	var sensor Sensor
	// Name can be found in temp*_label
	fileBaseName := strings.TrimSuffix(name, "_input")
	nameData, readErr := os.ReadFile(sensorDir + "/" + fileBaseName + "_label")
	if readErr != nil {
		return
	}
	sensor.Name = string(nameData[:len(nameData)-1]) // trim the newline
	sensor.File, readErr = os.OpenFile(sensorDir+"/"+name, os.O_RDONLY, 0600)
	if readErr != nil {
		log.Println("Error opening sensor file: ", readErr)
		return
	}
	*sensors = append(*sensors, sensor)
}

func findSensors(sensorDir string) []Sensor {
	var sensors []Sensor

	// Scan sensor directory
	dirEntry, err := os.ReadDir(sensorDir)
	if err != nil {
		panic(err)
	}

	// Find all the temp*_input files
	for _, dir := range dirEntry {
		// Skip directories that aren't temp*_input
		if !strings.HasPrefix(dir.Name(), "temp") && strings.HasSuffix(dir.Name(), "input") {
			continue
		}
		getSensor(sensorDir, dir.Name(), &sensors)
	}

	return sensors
}

func MonitorCPUTemp(packetMaps *map[string]*socket.Packet) {
	sensorDir, found := getSensorDir()
	if !found {
		log.Println("/!\\ No CPU temperature sensor found")
		return
	}

	sensors := findSensors(sensorDir)
	if len(sensors) == 0 {
		log.Println("/!\\ No CPU temperature sensor found in ", sensorDir)
		return
	}

	for {
		// Read all the temp*_input files and their corresponding temp*_label files
		for _, sensor := range sensors {
			temp := sensor.Read()
			packet := getCPUPacket(packetMaps, sensor.Name) // trim the newline
			packet.SetTemperature(temp / 1000)
		}
		time.Sleep(REFRESH_DELAY)
	}
}
