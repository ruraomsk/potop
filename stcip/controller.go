package stcip

import (
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/stcip/transport"
)

var live chan any
var state StateCentral
var mutex sync.Mutex

type StateCentral struct {
	connect         bool
	autonom         bool
	CommandAllRed   int
	CommandFlashing int
	CommandDark     int
	CommandPhase    int
	CommandPlane    int

	Rphase  RepPhase
	Rplan   RepPlan
	Rmajor  RepMajor
	Ralarm  RepAlarm
	Rsg     RepSignalGroups
	Rstatus RepStatus
	Rsource RepSource
}

func getDuration() time.Duration {
	return 60 * time.Second
}
func SetAutonom(set bool) {
	state.autonom = set
}
func GetAutonom() bool {
	return state.autonom
}
func GetConnect() bool {
	return state.connect
}
func controlCentral() {
	var timer *time.Timer
	hardware.SetWork <- 0
	logger.Info.Print("Нет управления от центра")
	journal.SendMessage(2, "Нет управления от STCIP")
	state.connect = false
	for {
		<-live
		hardware.SetWork <- 1
		timer = time.NewTimer(getDuration())
		logger.Info.Print("Есть управление от центра")
		journal.SendMessage(2, "Есть управление от STCIP")
		state.connect = true
	loop:
		for {
			select {
			case <-timer.C:
				if !state.autonom {
					hardware.SetWork <- 0
					break loop
				}
			case <-live:
				timer.Stop()
				timer = time.NewTimer(getDuration())
			}
		}
		state.connect = false
		logger.Error.Print("Потеряно управление от центра")
		journal.SendMessage(2, "Потеряно управление от центра")
	}

}
func executeCommand(command transport.Command) {
	live <- 0
	mutex.Lock()
	defer mutex.Unlock()

	switch command.Code {
	case transport.CodeCallAllRed:
		hardware.CommandToKDM(4, command.Value)
		state.CommandAllRed = command.Value
	case transport.CodeCallFlash:
		hardware.CommandToKDM(3, command.Value)
		state.CommandFlashing = command.Value
	case transport.CodeCallDark:
		hardware.CommandToKDM(6, command.Value)
		state.CommandDark = command.Value
	case transport.CodeCallPlan:
		hardware.CommandToKDM(7, command.Value)
		state.CommandPlane = command.Value
	case transport.CodeCallPhase:
		hardware.CommandToKDM(8, command.Value)
		state.CommandPhase = command.Value
	default:
		logger.Error.Printf("not command %v", command)
	}
}
func executeRequest(request transport.Request) {
	live <- 0
	mutex.Lock()
	defer mutex.Unlock()

	switch request.Code {
	case transport.CodeReqPhase:
		transport.Sender <- state.Rphase.send()
	case transport.CodeReqPlan:
		transport.Sender <- state.Rplan.send()
	case transport.CodeReqStatus:
		transport.Sender <- state.Rstatus.send()
	case transport.CodeReqSource:
		transport.Sender <- state.Rsource.send()
	case transport.CodeReqSignalGroups:
		transport.Sender <- state.Rsg.send()
	case transport.CodeReqAlarm:
		transport.Sender <- state.Ralarm.send()
	case transport.CodeReqMajor:
		transport.Sender <- state.Rmajor.send()
	default:
		logger.Error.Printf("not request %v", request)
	}

}
func Start() {
	go transport.Transport()
	live = make(chan any)
	go controlCentral()
	for {
		select {
		case command := <-transport.Commander:
			if !state.autonom {
				executeCommand(command)
			}
		case request := <-transport.Requester:
			if !state.autonom {
				executeRequest(request)
			}
		}
	}
}
