package watchdogs

import (
	"os"

	"github.com/bettercallmolly/illustrious/configuration"
	"github.com/bettercallmolly/illustrious/socket"
	"github.com/fsnotify/fsnotify"
)

// check if file '/run/systemd/units/invocation:{name}.service' exists

func setServiceState(serviceName string, packet *socket.Packet) {
	_, err := os.Lstat("/run/systemd/units/invocation:" + serviceName + ".service")
	if err != nil {
		packet.SetState(socket.S_DEAD)
	} else {
		packet.SetState(socket.S_OK)
	}
}

func ManualServices(packetMaps *map[string]*socket.Packet) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	for _, service := range configuration.LoadedConfiguration.MonitoredServices {
		serviceSocket := socket.NewMonitoringPacket(socket.C_SERVICE, socket.DT_UINT8, service)
		(*packetMaps)[service] = serviceSocket
		// Since we're iterating over files, there cannot be duplicates, `watcher.Add` cannot fail, so we can ignore the error
		watcher.Add("/run/systemd/units")
		// Check if service's file exists
		setServiceState(service, serviceSocket)
	}

	go func() {
		for {
			select {
			case changes := <-configuration.ServicesChanges:
				for serviceName, added := range changes {
					if added {
						serviceSocket := socket.NewMonitoringPacket(socket.C_SERVICE, socket.DT_UINT8, serviceName)
						(*packetMaps)[serviceName] = serviceSocket
						watcher.Add("/run/systemd/units")
						setServiceState(serviceName, serviceSocket)
						socket.ConnectedClients.Broadcast(serviceSocket.GetRawBytes())
					} else {
						packet, ok := (*packetMaps)[serviceName]
						if !ok {
							continue
						}
						packet.RemoveDOM()
						delete(*packetMaps, serviceName)
					}
				}
			}
		}
	}()

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
				socket.ConnectedClients.Broadcast(packet.GetRawBytes())
			} else if event.Op.Has(fsnotify.Remove) {
				(*packetMaps)[serviceName].SetState(socket.S_DEAD)
				packet := (*packetMaps)[serviceName]
				socket.ConnectedClients.Broadcast(packet.GetRawBytes())
			}
		}
	}
}
