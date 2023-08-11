package watchdogs

import (
	"context"
	"log"
	"strings"

	"github.com/bettercallmolly/illustrious/socket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func getDockerPacket(packetMaps *map[string]*socket.Packet, containerName string) *socket.Packet {
	containerName = strings.TrimPrefix(containerName, "/")
	packet, ok := (*packetMaps)[containerName]
	// If the packet does not exist, create it
	if !ok {
		(*packetMaps)[containerName] = socket.NewPacket(socket.C_SERVICE, socket.DT_UINT8, containerName)
		return (*packetMaps)[containerName]
	} else {
		return packet
	}
}

func MonitorContainers(packetMaps *map[string]*socket.Packet, clients *socket.Clients) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// Register packets for every container that already exists
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		log.Println("/!\\ Could not list containers", err)
	}
	for _, container := range containers {
		packet := getDockerPacket(packetMaps, container.Names[0])
		inspection, err := client.ContainerInspect(ctx, container.ID)
		if err != nil {
			log.Println("/!\\ Could not inspect container", container.Names[0])
			continue
		}
		// Check if the container has a healthcheck, and if so, report 'healthy' or 'unhealthy'
		if inspection.State.Health != nil {
			if inspection.State.Health.Status == "healthy" {
				packet.SetState(socket.S_HEALTHY)
			} else {
				packet.SetState(socket.S_UNHEALTHY)
			}
			continue
		}
		// If the container does not have a healthcheck, report 'running' or 'stopped'
		if inspection.State.Running {
			packet.SetState(socket.S_HEALTHY)
		} else {
			packet.SetState(socket.S_DEAD)
		}
	}
	// Listen to events 'destroy', 'die', 'health_staus', 'start'
	eventFilters := filters.NewArgs()
	eventFilters.Add("type", "container")
	eventFilters.Add("event", "destroy")
	eventFilters.Add("event", "die")
	eventFilters.Add("event", "health_status")
	eventFilters.Add("event", "start")
	eventChan, errChan := client.Events(ctx, types.EventsOptions{
		Filters: eventFilters,
	})
	for { // Poll events
		select {
		case event := <-eventChan:
			if event.Type == "container" {
				// Get the packet corresponding to the container
				packet := getDockerPacket(packetMaps, event.Actor.Attributes["name"])
				if event.Action == "health_status" {
					if event.Actor.Attributes["health_status"] == "healthy" {
						packet.SetState(socket.S_HEALTHY)
					} else {
						packet.SetState(socket.S_UNHEALTHY)
					}
				} else if event.Action == "start" {
					packet.SetState(socket.S_HEALTHY)
				} else if event.Action == "die" {
					packet.SetState(socket.S_DEAD)
				} else if event.Action == "destroy" {
					packet.SetState(socket.S_DEAD)
				}
				clients.Broadcast(packet.GetRawBytes())
			}
		case err := <-errChan:
			if err != nil {
				log.Println("/!\\ Could not listen to events", err)
			}
		}
	}
}
