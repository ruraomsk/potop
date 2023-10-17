package stat

import "github.com/ruraomsk/ag-server/logger"

var InStat chan OneTick
var Statistics Chanels

func Start(chanels int) {
	InStat = make(chan OneTick, 100)
	Statistics.clear(chanels)
	for {
		select {
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
	Statistics.clear(0)
	for {
		select {
		case t := <-InStat:
			logger.Error.Printf("статистика отключена %v", t)
		}
	}
}
