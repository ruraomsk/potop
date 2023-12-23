package stcip

import (
	"fmt"
	"time"
)

func timeToString(t time.Time) string {
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func GetState() StateCentral {
	mutex.Lock()
	defer mutex.Unlock()
	return state
}
