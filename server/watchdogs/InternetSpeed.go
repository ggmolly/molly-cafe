package watchdogs

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorInternetSpeed(packet *socket.Packet) {
	for {
		resp, err := http.Get(os.Getenv("SPEEDTEST_URL"))
		if err != nil {
			log.Println("/!\\ Error while opening HTTP request: ", err)
		} else {
			null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
			start := time.Now()
			written, _ := io.Copy(null, resp.Body)
			elapsed := time.Since(start)
			// Calculate the speed in Mbps
			speed := float64(written) / elapsed.Seconds() / 125000
			packet.SetUint32(uint32(speed))
			socket.ConnectedClients.Broadcast(packet.GetRawBytes())
			null.Close()
		}
		time.Sleep(1 * time.Hour) // Let's not rekt their bandwidth for no reason :)
	}
}
