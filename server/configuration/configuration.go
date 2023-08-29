package configuration

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

var (
	LoadedConfiguration = Configuration{}

	oldMonitoredServices = []string{}

	// Will be used to notify the services watchdogs that the configuration has changed
	// The map contains the services that have been added or removed (true if added, false if removed)
	ServicesChanges = make(chan map[string]bool)
)

type Configuration struct {
	MonitoredServices []string          `json:"services"`
	DiskTranslations  map[string]string `json:"disk_translations"`
}

func contains(slice *[]string, element string) bool {
	for _, item := range *slice {
		if item == element {
			return true
		}
	}
	return false
}

func LoadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Println("Failed to read the config.json file")
		return
	}
	err = json.Unmarshal(data, &LoadedConfiguration)
	if err != nil {
		log.Println("Failed to parse the config.json file")
		return
	}
	log.Println("Configuration loaded !")
	log.Println("Monitored services:")
	for _, service := range LoadedConfiguration.MonitoredServices {
		log.Println("  -", service)
	}
	var change = make(map[string]bool)
	for _, service := range LoadedConfiguration.MonitoredServices {
		if !contains(&oldMonitoredServices, service) {
			change[service] = true
		}
	}
	for _, service := range oldMonitoredServices {
		if !contains(&LoadedConfiguration.MonitoredServices, service) {
			change[service] = false
		}
	}
	if len(change) > 0 {
		// We're using a goroutine to avoid blocking the main thread
		// the channel will not be read until Configuration package inits
		// so we don't have to worry about the channel being full
		go func() {
			ServicesChanges <- change
		}()
	}
	log.Println("Disk translations:")
	for key, value := range LoadedConfiguration.DiskTranslations {
		log.Println("  -", key, "->", value)
	}
	oldMonitoredServices = LoadedConfiguration.MonitoredServices
}

func PollConfig() {
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
				log.Println("Config file changed, reloading...")
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
	LoadConfig()
	go PollConfig()
}
