package transport

import (
	"time"
)

func senderReplay() {
	for {
		s := <-Sender
		toHistory <- History{Time: time.Now(), Type: 2, Message: s.toString()}
	}
}
