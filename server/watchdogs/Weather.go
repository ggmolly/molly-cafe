package watchdogs

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/bettercallmolly/illustrious/socket"
)

const (
	MAX_WIND_SPEED = 30 // m/s (108 km/h)
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
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

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

func getWeatherPacket(packetMaps *map[string]*socket.Packet) *socket.Packet {
	packet, ok := (*packetMaps)["weather"]
	if !ok {
		packet = socket.NewUntrackedPacket(socket.T_WEATHER, 0x00, socket.DT_SPECIAL, "")
		(*packetMaps)["weather"] = packet
	}
	return packet
}

func getWeatherData() (OWM_Data, error) {
	// Do the request to OpenWeatherMap
	// var data OWM_Data
	// resp, err := http.Get("https://samples.openweathermap.org/data/2.5/find?q=London&units=metric&appid=b1b15e88fa797225412429c1c50c122a1")
	// if err != nil {
	// 	return data, err
	// }
	// defer resp.Body.Close()
	// // Serialize the response
	// err = json.NewDecoder(resp.Body).Decode(&data)
	// return data, err
	file, _ := os.OpenFile("./test_data.json", os.O_RDONLY, 0644)
	defer file.Close()
	var data OWM_Data
	err := json.NewDecoder(file).Decode(&data)
	return data, err
}

func serializeWeatherPacket(buffer *bytes.Buffer) error {
	data, err := getWeatherData()
	if err != nil {
		return err
	}
	// Normalize the wind speed
	windSpeed := uint8(data.Wind.Speed)
	if windSpeed > MAX_WIND_SPEED {
		windSpeed = 255
	} else {
		windSpeed = uint8(windSpeed / MAX_WIND_SPEED * 255)
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

	// Serialize the length of the weather condition (1 byte)
	buffer.WriteByte(uint8(len(data.Weather[0].Description)))

	// Serialize the weather condition
	buffer.WriteString(data.Weather[0].Description)
	return nil
}

func MonitorWeather(packetMaps *map[string]*socket.Packet) {
	// TODO: Get real data from OpenWeatherMap, use a caching file to avoid spamming the API (one call / hour)
	// use gob encoding to store old file along with the request timestamp
	// XXX: use https://samples.openweathermap.org/data/2.5/find?q=London&units=metric&appid=b1b15e88fa797225412429c1c50c122a1&units=metric

	// Get the packet from the map, or create it if it doesn't exist
	packet := getWeatherPacket(packetMaps)
	var buffer bytes.Buffer
	if err := serializeWeatherPacket(&buffer); err != nil {
		log.Println("Failed to serialize weather packet:", err)
		return
	}
	packet.Data = buffer.Bytes()
	packet.Dirty = true
}
