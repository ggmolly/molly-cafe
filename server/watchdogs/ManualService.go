package watchdogs

import (
	"os"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/fsnotify/fsnotify"
)

// check if file '/run/systemd/units/invocation:{name}.service' exists

func ManualServices(packetMaps *map[string]*socket.Packet, clients *socket.Clients, services ...string) {
	watcher, err := fsnotify.NewWatcher()
	for _, service := range services {
		serviceSocket := socket.NewPacket(socket.C_SERVICE, socket.DT_UINT8, service)
		(*packetMaps)[service] = serviceSocket
		if err != nil {
			panic(err)
		}
		err = watcher.Add("/run/systemd/units")
		if err != nil {
			panic(err)
		}
		// Check if service's file exists
		_, err := os.Lstat("/run/systemd/units/invocation:" + service + ".service")
		if err != nil {
			serviceSocket.SetState(socket.S_DEAD)
		} else {
			serviceSocket.SetState(socket.S_OK)
		}
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if len(event.Name) < 38 {
				continue
			}
			serviceName := event.Name[30 : len(event.Name)-8]
			// Check if service is in the map
			_, ok = (*packetMaps)[serviceName]
			if !ok {
				continue
			}
			if event.Op.Has(fsnotify.Create) {
				(*packetMaps)[serviceName].SetState(socket.S_OK)
				// Broadcast to clients
				packet := (*packetMaps)[serviceName]
				clients.Broadcast(packet.GetRawBytes())
			} else if event.Op.Has(fsnotify.Remove) {
				(*packetMaps)[serviceName].SetState(socket.S_DEAD)
				packet := (*packetMaps)[serviceName]
				clients.Broadcast(packet.GetRawBytes())
			}
		}
	}
}
