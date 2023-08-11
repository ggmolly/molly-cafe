package watchdogs

import (
	"log"
	"os/exec"
	"time"

	"github.com/bettercallmolly/illustrious/socket"
)

func MonitorLoggedUsers(packet *socket.Packet) {
	for {
		// FIXME: stop using who please

		cmd, err := exec.Command("who").Output()
		if err != nil {
			log.Println("/!\\ Error while executing who command")
		}

		// Count lines
		lines := 0
		for _, c := range cmd {
			if c == '\n' {
				lines++
			}
		}

		packet.SetUint32(uint32(lines))
		time.Sleep(REFRESH_DELAY)
	}
}
