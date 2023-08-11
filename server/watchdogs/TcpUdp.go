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
	for {
		tcp, err := os.OpenFile("/proc/net/tcp", os.O_RDONLY, 0)
		if err != nil {
			log.Fatal("/!\\ Error while opening /proc/net/tcp", err)
		}
		udp, err := os.OpenFile("/proc/net/udp", os.O_RDONLY, 0)
		if err != nil {
			log.Fatal("/!\\ Error while opening /proc/net/tcp", err)
		}
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
		tcp.Close()
		udp.Close()
		time.Sleep(REFRESH_DELAY)
	}
}
