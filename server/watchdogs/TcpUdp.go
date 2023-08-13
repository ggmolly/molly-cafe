package watchdogs

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorSocketConnections(tcpPacket, udpPacket *socket.Packet) {
	log.Println("Monitoring socket connections...")
	tcp, tcpErr := os.OpenFile("/proc/net/tcp", os.O_RDONLY, 0)
	udp, udpErr := os.OpenFile("/proc/net/udp", os.O_RDONLY, 0)
	if tcpErr != nil {
		log.Fatal("/!\\ Error while opening /proc/net/tcp", tcpErr)
	}
	if udpErr != nil {
		log.Fatal("/!\\ Error while opening /proc/net/udp", udpErr)
	}
	defer tcp.Close()
	defer udp.Close()
	for {
		tcp.Seek(0, 0)
		udp.Seek(0, 0)
		// Count number of lines in each files
		var tcpCount, udpCount uint32
		for {
			var buffer [1024]byte
			n, err := tcp.Read(buffer[:])
			if err != nil {
				break
			}
			tcpCount += uint32(bytes.Count(buffer[:n], []byte{'\n'}))
		}
		for {
			var buffer [1024]byte
			n, err := udp.Read(buffer[:])
			if err != nil {
				break
			}
			udpCount += uint32(bytes.Count(buffer[:n], []byte{'\n'}))
		}
		tcpPacket.SetUint32(tcpCount)
		udpPacket.SetUint32(udpCount)
		time.Sleep(REFRESH_DELAY)
	}
}
