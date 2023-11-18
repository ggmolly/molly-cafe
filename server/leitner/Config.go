package leitner

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

const (
	configPath = "leitner.json"
)

var (
	LeitnerData LeitnerConfig
	configMutex = &sync.Mutex{}
)

type LeitnerTopic struct {
	CompletedCards uint32 `json:"completed_cards"`
	Total          uint32 `json:"total"`
}

type LeitnerConfig struct {
	Topics    map[string]LeitnerTopic `json:"topics"`
	Streak    uint32                  `json:"streak"`
	UpdatedAt time.Time               `json:"updated_at"`
}

func (c *LeitnerConfig) Save() error {
	configMutex.Lock()
	defer configMutex.Unlock()
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(c)
}

func (c *LeitnerConfig) Load() error {
	configMutex.Lock()
	defer configMutex.Unlock()
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(c)
}

// Returns the topic and whether it exists
func (c *LeitnerConfig) GetTopic(name string) (LeitnerTopic, bool) {
	configMutex.Lock()
	defer configMutex.Unlock()
	topic, ok := c.Topics[name]
	return topic, ok
}

// Sets a topic
func (c *LeitnerConfig) SetTopic(name string, topic LeitnerTopic) {
	configMutex.Lock()
	defer configMutex.Unlock()
	c.Topics[name] = topic
	c.Save() // XXX: This might cause an infinite loop because of the lock in the Save() function
}

// Get last update time
func (c *LeitnerConfig) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}

// Refresh update time
func (c *LeitnerConfig) RefreshUpdatedAt() {
	configMutex.Lock()
	defer configMutex.Unlock()
	c.UpdatedAt = time.Now()
	// how stupid is this ?
	c.UpdatedAt = c.UpdatedAt.Add(-time.Duration(c.UpdatedAt.Hour()) * time.Hour)
	c.UpdatedAt = c.UpdatedAt.Add(-time.Duration(c.UpdatedAt.Minute()) * time.Minute)
	c.UpdatedAt = c.UpdatedAt.Add(-time.Duration(c.UpdatedAt.Second()) * time.Second)
	c.UpdatedAt = c.UpdatedAt.Add(-time.Duration(c.UpdatedAt.Nanosecond()) * time.Nanosecond)
}

func (c *LeitnerConfig) GetStreak() uint32 {
	return c.Streak
}

func init() {
	// check if the file leitner.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Leitner config does not exist. The Leitner API will not be available.")
	} else {
		// Load the config
		if err := LeitnerData.Load(); err != nil {
			log.Println("Could not load leitner config. The Leitner API will not be available.")
		}
	}
	// Update the packets
	for topic := range LeitnerData.Topics {
		UpdateLeitnerPacket(topic)
	}
	UpdateStreakPacket(LeitnerData.Streak)
}
