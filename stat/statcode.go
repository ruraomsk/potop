package stat

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

var InStat chan OneTick
var Statistics = Chanels{chanels: make(map[int]*OneChanel)}

func Start(chanels int, diaps int) {
	InStat = make(chan OneTick, 100)
	Statistics.clearAll(chanels, diaps)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			Statistics.newSecond()
		case t := <-InStat:
			err := Statistics.add(t)
			if err != nil {
				logger.Error.Print(err.Error())
			}
		}
	}
}
func NoStatistics() {
	InStat = make(chan OneTick, 100)
	Statistics.clearAll(0, 0)
	for {
		select {
		case t := <-InStat:
			logger.Error.Printf("статистика отключена %v", t)
		}
	}
}
