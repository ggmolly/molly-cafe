package configuration

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	LoadedConfiguration = Configuration{}
	// Will be used to notify the services watchdogs that the configuration has changed
	// The map contains the services that have been added or removed (true if added, false if removed)
	ServicesChanges = make(chan map[string]bool)
)

type Configuration struct {
	MonitoredServices []string          `json:"services"`
	DiskTranslations  map[string]string `json:"disk_translations"`
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func LoadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return
	}
	// Deep copy the old services
	var oldServices = make([]string, len(LoadedConfiguration.MonitoredServices))
	copy(oldServices, LoadedConfiguration.MonitoredServices)
	err = json.Unmarshal(data, &LoadedConfiguration)
	if err != nil {
		return
	}
	var changes = make(map[string]bool)
	for _, service := range LoadedConfiguration.MonitoredServices {
		if !contains(oldServices, service) {
			changes[service] = true
		}
	}
	for _, service := range oldServices {
		if !contains(LoadedConfiguration.MonitoredServices, service) {
			changes[service] = false
		}
	}
	if len(changes) > 0 {
		// We're using a goroutine to avoid blocking the main thread
		// the channel will not be read until Configuration package inits
		// so we don't have to worry about the channel being full
		go func() {
			ServicesChanges <- changes
		}()
	}
}

func PollConfig() {
	LoadConfig()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Failed to create a watcher for the config file")
		return
	}
	defer watcher.Close()
	err = watcher.Add("config.json")
	if err != nil {
		log.Println("Failed to add the config file to the watcher")
		return
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Failed to get the event from the watcher")
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				LoadConfig()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Println("Failed to get the error from the watcher")
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func init() {
	log.Println("Loading configuration...")
	go PollConfig()
}
