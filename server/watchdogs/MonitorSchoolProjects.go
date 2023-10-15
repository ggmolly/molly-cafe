package watchdogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/fsnotify/fsnotify"
)

type SchoolProject struct {
	Name        string `json:"name"`
	Href        string `json:"href"`
	Description string `json:"description"`
	Wip         bool   `json:"wip"`
	Grading     bool   `json:"grading"`
	Grade       int8   `json:"grade"`
}

// Create / update the packet in the packet map and sends it to the clients
func updateProject(packetMap *socket.T_PacketMap, project SchoolProject, projectPath string) {
	packet, ok := packetMap.GetPacketByName(projectPath)
	if !ok {
		packet = socket.NewPacket(socket.T_SCHOOL_PROJECTS, socket.C_SCHOOL, socket.DT_SPECIAL, project.Name)
	}
	var buffer bytes.Buffer

	// int8 -> wip / grading (bitmask)
	var mask uint8 = 0x00
	if project.Wip {
		mask |= 0x01
	}
	if project.Grading {
		mask |= 0x02
	}
	buffer.WriteByte(mask)

	// int8 -> grade (whatever if not grading)
	buffer.WriteByte(uint8(project.Grade))

	// uint8 -> description length
	buffer.WriteByte(uint8(len(project.Description)))

	// string -> description
	buffer.WriteString(project.Description)

	// uint8 -> href length
	buffer.WriteByte(uint8(len(project.Href)))

	// string -> href
	buffer.WriteString(project.Href)

	packet.Data = buffer.Bytes()
	packet.Name = project.Name
	packetMap.AddPacket(projectPath, packet)

	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func readProject(path string) (SchoolProject, error) {
	var project SchoolProject
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return project, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	for {
		if err := decoder.Decode(&project); err == io.EOF {
			break
		} else if err != nil {
			return project, err
		}
	}
	if project.Name == "" {
		return project, fmt.Errorf("project name is empty")
	}
	return project, nil
}

func MonitorSchoolProjects(packetMap *socket.T_PacketMap, rootPath string) {
	// Find all JSON files in the projects folder and add them to the packet map
	files, err := filepath.Glob(filepath.Join(rootPath, "*.json"))
	if err != nil {
		log.Println("Failed to find projects in the projects folder")
		return
	}
	for _, file := range files {
		if project, err := readProject(file); err != nil {
			log.Println("Failed to read project", file)
		} else {
			updateProject(packetMap, project, file)
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to create a new watcher for the projects folder")
		return
	}
	defer watcher.Close()

	err = watcher.Add(rootPath)
	if err != nil {
		log.Println("Failed to add the projects folder to the watcher")
		return
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op.Has(fsnotify.Create) || event.Op.Has(fsnotify.Write) {
				if project, err := readProject(event.Name); err == nil {
					updateProject(packetMap, project, event.Name)
				}
			}
			if event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) {
				if packet, ok := packetMap.GetPacketByName(event.Name); ok {
					packet.RemoveDOM()
					delete(*packetMap, event.Name)
				}
			}
		}
	}
}
