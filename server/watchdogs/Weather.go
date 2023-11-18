package watchdogs

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bettercallmolly/illustrious/configuration"
	"github.com/bettercallmolly/illustrious/socket"
)

const (
	WEATHER_CACHE_DURATION = 5 * time.Minute
)

type OWM_Data struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type WeatherCache struct {
	Data       OWM_Data
	LastUpdate time.Time
}

var (
	CachedWeatherData WeatherCache
	lastRainedAt      time.Time
	requestURL        string
)

// Very naive function to compute the rain intensity (who cares)
func computeRainIntensity(data OWM_Data) uint8 {
	// https://openweathermap.org/weather-conditions

	weatherId := data.Weather[0].ID
	weatherDesc := data.Weather[0].Description

	// 1. Check if the weather ID corresponds to a rainy weather
	if !(weatherId >= 200 || weatherId < 700) {
		return 0
	}

	containsDrizzle := strings.Contains(weatherDesc, "drizzle")
	containsRain := strings.Contains(weatherDesc, "rain")
	if !containsDrizzle && !containsRain {
		return 0
	}

	var intensity uint8 = 0
	if strings.HasPrefix(weatherDesc, "extreme") {
		intensity = 6
	}
	if strings.HasPrefix(weatherDesc, "very heavy") {
		intensity = 5
	}
	if strings.HasPrefix(weatherDesc, "heavy intensity") {
		intensity = 4
	}
	if strings.HasPrefix(weatherDesc, "moderate") {
		intensity = 3
	}
	if strings.Contains(weatherDesc, "light") {
		intensity = 1
	}
	if intensity == 0 {
		intensity = 2 // no specific adjective
	}

	var score uint8 = 0
	if containsDrizzle {
		score = 20
	}
	if containsRain {
		score = 42
	}
	// Max with drizzle : 20 * 6 = 120
	// Max with rain : 42 * 6 = 252
	return intensity * score
}

func getWeatherPacket(packetMaps *socket.T_PacketMap) *socket.Packet {
	packet, ok := packetMaps.GetPacketByName("weather")
	if !ok {
		packet = socket.NewUntrackedPacket(socket.T_WEATHER, 0x00, socket.DT_SPECIAL, "")
		packetMaps.AddPacket("weather", packet)
	}
	return packet
}

func getWeatherData() (OWM_Data, error) {
	// Check if the data is cached
	if time.Since(CachedWeatherData.LastUpdate) < WEATHER_CACHE_DURATION {
		return CachedWeatherData.Data, nil
	}
	// Do the request to OpenWeatherMap
	var data OWM_Data
	resp, err := http.Get(requestURL)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	// Serialize the response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Println("Failed to decode weather response:", err)
		return data, err
	}
	// Save the data in the cache
	CachedWeatherData.Data = data
	CachedWeatherData.LastUpdate = time.Now()
	// Save the cache on disk
	file, err := os.OpenFile(".cached_weather", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed to open weather cache:", err)
		return data, err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(CachedWeatherData)
	if err != nil {
		log.Println("Failed to encode weather cache:", err)
	}
	return data, err
}

func serializeWeatherPacket(buffer *bytes.Buffer) error {
	data, err := getWeatherData()
	if err != nil {
		return err
	}
	// Normalize the wind speed
	windSpeed := uint8(data.Wind.Speed)
	if data.Wind.Speed > 255 {
		windSpeed = 255
	}
	buffer.WriteByte(windSpeed)

	// Compute rain intensity
	rainIntensity := computeRainIntensity(data)
	buffer.WriteByte(rainIntensity)

	// Serialize the cloudiness
	buffer.WriteByte(uint8(data.Clouds.All))

	// Serialize the temperature as two uint8 (first byte = real part, second byte = decimal part)
	buffer.WriteByte(uint8(data.Main.Temp))
	buffer.WriteByte(uint8(data.Main.Temp*100) % 100)

	// Serialize the felt temperature (same process)
	buffer.WriteByte(uint8(data.Main.FeelsLike))
	buffer.WriteByte(uint8(data.Main.FeelsLike*100) % 100)

	// Serialize the humidity (same process)
	buffer.WriteByte(uint8(data.Main.Humidity))
	buffer.WriteByte(uint8(data.Main.Humidity*100) % 100)

	// Serialize the time to sunrise (4 bytes)
	buffer.WriteByte(uint8(data.Sys.Sunrise >> 24))
	buffer.WriteByte(uint8(data.Sys.Sunrise >> 16))
	buffer.WriteByte(uint8(data.Sys.Sunrise >> 8))
	buffer.WriteByte(uint8(data.Sys.Sunrise))

	// Serialize the time to sunset (4 bytes)
	buffer.WriteByte(uint8(data.Sys.Sunset >> 24))
	buffer.WriteByte(uint8(data.Sys.Sunset >> 16))
	buffer.WriteByte(uint8(data.Sys.Sunset >> 8))
	buffer.WriteByte(uint8(data.Sys.Sunset))

	// Serialize the current date (4 bytes)
	buffer.WriteByte(uint8(data.Dt >> 24))
	buffer.WriteByte(uint8(data.Dt >> 16))
	buffer.WriteByte(uint8(data.Dt >> 8))
	buffer.WriteByte(uint8(data.Dt))

	// Serialize the last rain time (2 bytes)
	var lastRainTime int16
	// If it's been more than 30m, set to -1
	if time.Since(lastRainedAt).Seconds() > 1800 {
		lastRainTime = -1
	} else {
		lastRainTime = int16(time.Since(lastRainedAt).Seconds())
	}
	// If it's currently raining, set to 0 and update the last rain time
	if rainIntensity > 0 {
		lastRainTime = 0
		lastRainedAt = time.Now()
	}
	buffer.WriteByte(uint8(lastRainTime >> 8))
	buffer.WriteByte(uint8(lastRainTime))

	// Serialize the length of the weather condition (1 byte)
	buffer.WriteByte(uint8(len(data.Weather[0].Description)))

	// Serialize the weather condition
	buffer.WriteString(data.Weather[0].Description)
	return nil
}

func MonitorWeather(packetMaps *socket.T_PacketMap) {
	// Get the packet from the map, or create it if it doesn't exist
	for {
		packet := getWeatherPacket(packetMaps)
		var buffer bytes.Buffer
		if err := serializeWeatherPacket(&buffer); err != nil {
			log.Println("Failed to serialize weather packet:", err)
			return
		}
		packet.Data = buffer.Bytes()
		packet.Dirty = true
		// Sleep until the next update with a 1 second margin
		time.Sleep(WEATHER_CACHE_DURATION - time.Since(CachedWeatherData.LastUpdate) + 1*time.Second)
	}
}

func init() {
	requestURL = fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%.6f&lon=%.6f&appid=%s&units=metric&lang=en",
		configuration.LoadedConfiguration.OpenWeatherMap.Latitude,
		configuration.LoadedConfiguration.OpenWeatherMap.Longitude,
		configuration.LoadedConfiguration.OpenWeatherMap.API,
	)
	file, err := os.OpenFile(".cached_weather", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Failed to open weather cache:", err)
		return
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&CachedWeatherData)
	if err != nil {
		log.Println("Failed to decode weather cache:", err)
		return
	}
	if computeRainIntensity(CachedWeatherData.Data) > 0 {
		lastRainedAt = CachedWeatherData.LastUpdate
	}
	log.Println("Weather cache loaded! Cached at:", CachedWeatherData.LastUpdate)
}
