package watchdogs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/fsnotify/fsnotify"
)

var (
	titleRegex = regexp.MustCompile("<title>(.+)</title>")
)

type PistachePost struct {
	Title        string
	CreationDate time.Time
}

// Reads first bytes of the HTML pistache post to get the <title> tag
func getTitle(path string) (string, error) {
	buffer := make([]byte, 1024)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	matches := titleRegex.FindSubmatch(buffer)
	if len(matches) < 2 {
		return "", fmt.Errorf("no title found in %s", path)
	}
	return string(matches[1]), nil
}

func updatePost(packetMap *socket.T_PacketMap, path, title string) {
	packet, ok := packetMap.GetPacketByName(path)
	if !ok {
		packet = socket.NewPacket(socket.T_PISTACHE, socket.C_PISTACHE, socket.DT_SPECIAL, path)
	}
	var buffer bytes.Buffer

	// uint16 -> href length
	name := filepath.Base(path)
	length := uint16(len(name))
	buffer.WriteByte(byte(length >> 8))
	buffer.WriteByte(byte(length))

	// string -> href
	buffer.WriteString(name)

	// UTC timestamp (uint32) -> creation date
	var timestamp uint32
	stat, err := os.Stat(path)
	if err == nil {
		timestamp = uint32(stat.ModTime().Unix())
	}
	buffer.WriteByte(byte(timestamp >> 24))
	buffer.WriteByte(byte(timestamp >> 16))
	buffer.WriteByte(byte(timestamp >> 8))
	buffer.WriteByte(byte(timestamp))
	packet.Data = buffer.Bytes()
	packet.Name = title
	packetMap.AddPacket(path, packet)
	socket.ConnectedClients.Broadcast(packet.GetRawBytes())
}

func MonitorPistachePosts(packetMap *socket.T_PacketMap, rootPath string) {
	files, err := filepath.Glob(filepath.Join(rootPath, "*.html"))
	if err != nil {
		return
	}
	for _, file := range files {
		if title, err := getTitle(file); err != nil {
			continue
		} else {
			updatePost(packetMap, file, title)
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	err = watcher.Add(rootPath)
	if err != nil {
		return
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op.Has(fsnotify.Create) || event.Op.Has(fsnotify.Write) {
				if title, err := getTitle(event.Name); err != nil {
					continue
				} else {
					updatePost(packetMap, event.Name, title)
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
