package utopia

import (
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/setup"
)

// From Spot to the controller					Reply from controller
// __________________________________________________________________________________
// TLC and group control: message 2				Status and detections: message 190
// (every second)
// ___________________________________________________________________________________
// Signal group count-down: message 8			Signal group feedback: message 4
// (every second)
// Extended count-down: message 9
// (every second)
// ___________________________________________________________________________________
// Diagnostic request: message 0				Extended diagnostic: message 5
// (event driven)
// ___________________________________________________________________________________
// Classified counts / vehicle length data		Classified counts / vehicle length data:
// request: message 24							message 24
// (periodic, default = 1 minute)
// ___________________________________________________________________________________
// Classified counts / vehicle speed data		Classified counts / vehicle length data:
// request: message 25							message 25
// (periodic, default = 1 minute)
// ___________________________________________________________________________________
// Special commands: message 6					Reply to a command: message 7
// ___________________________________________________________________________________
// Bus prediction: message 23					Bus detection: message 1
// (event driven)								or
// Time setting: message 3						Empty message (acknowledge)
// (periodic, default = 5 minutes)
// Empty message (polling)
// ___________________________________________________________________________________

var ctrl = ControllerUtopia{id: 1, lastACK: 0, input: make([]byte, 0), output: make([]byte, 0)}
var mutex sync.Mutex

func getDuration() time.Duration {
	return 20 * time.Second
}

var live chan any

func GetControllerUtopia() ControllerUtopia {
	mutex.Lock()
	defer mutex.Unlock()
	return ctrl
}

func controlUtopiaServer() {
	var timer *time.Timer
	hardware.SetWork <- 0
	logger.Info.Print("Нет управления от utopia")
	for {
		<-live
		hardware.SetWork <- 1
		timer = time.NewTimer(getDuration())
		logger.Info.Print("Есть управление от utopia")
		hardware.SetWork <- 1
	loop:
		for {
			select {
			case <-timer.C:
				hardware.SetWork <- 0
				break loop
			case <-live:
				timer.Stop()
				timer = time.NewTimer(getDuration())
			}
		}
		logger.Error.Print("Потеряно управление от utopia")
	}

}
func Controller() {
	live = make(chan any)
	go controlUtopiaServer()
	for {
		ctrl.input = <-fromServer
		live <- 0
		if hardware.IsConnectedKDM() {
			workMessage()
		} else {
			if setup.Set.Utopia.Replay {
				workMessage()
			}
		}
	}
}
func workMessage() {
	mutex.Lock()
	defer mutex.Unlock()
	if err := ctrl.verify(); err != nil {
		ctrl.sendNACK()
		logger.Error.Print(err.Error())
		return
	}
	if ctrl.isNak() {
		//Повторим предыдущее сообшение
		logger.Error.Printf("Повторяем сообщение % 02X", ctrl.output)
		toServer <- ctrl.output
		return
	}
	if ctrl.input[4] == ctrl.lastACK {
		logger.Error.Printf("ACK не изменился")
		ctrl.sendLive()
		return
	}
	if ctrl.isLive() {
		ctrl.sendLive()
		return
	}
	if !hardware.IsConnectedKDM() {
		ctrl.sendLive()
		return
	}
	switch ctrl.input[6] {
	case 2:
		// Message 2 – TLC and group control
		err := ctrl.TlcAndGroupControl.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.TlcAndGroupControl.execute()
		ctrl.status = ctrl.TlcAndGroupControl.command
		ctrl.StatusAndDetections.fill()
		ctrl.sendReplay(ctrl.StatusAndDetections.toData())
	case 8:
		// Message 8 – Signal group count-down
		err := ctrl.CountDown.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		if ctrl.status == 2 {
			ctrl.CountDown.execute()
		}
		ctrl.SignalGroupFeedback.fill()
		ctrl.sendReplay(ctrl.SignalGroupFeedback.toData())
	case 9:
		// Message 9 – Signal group extended count-down
		err := ctrl.ExtendedCountDown.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		if ctrl.status == 2 {
			ctrl.ExtendedCountDown.execute()
		}
		ctrl.SignalGroupFeedback.fill()
		ctrl.sendReplay(ctrl.SignalGroupFeedback.toData())
	case 3:
		// Message 3 – Date and time setting
		err := ctrl.DateAndTime.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		// fmt.Println(ctrl.DateAndTime.DateTime.String())
		ctrl.sendLive()
	case 0:
		// Message 0 – Diagnostic request message
		err := ctrl.DiagnosticRequest.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.ExtendedDiagnostic.fill()
		ctrl.sendReplay(ctrl.ExtendedDiagnostic.toData())
	case 24:
		// Message 24 – Request for classified counts by vehicle length
		err := ctrl.ReqClassifiedLenght.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.ClassifiedCounts.fill()
		ctrl.sendReplay(ctrl.ClassifiedCounts.toData())
	case 25:
		// Message 25 – Request for classified counts by vehicle speed
		err := ctrl.ReqClassifiedSpeed.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.ClassifiedSpeeds.fill()
		ctrl.sendReplay(ctrl.ClassifiedSpeeds.toData())
	case 6:
		// Message 6 - Special commands
		err := ctrl.SpecialCommands.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.ReplaySpecial.fill()
		ctrl.sendReplay(ctrl.ReplaySpecial.toData())
	case 23:
		// Message 23 – Bus prediction
		err := ctrl.BusPrediction.fromData(ctrl.data)
		if err != nil {
			logger.Error.Print(err.Error())
		}
		ctrl.BusDetection.fill()
		ctrl.sendReplay(ctrl.BusDetection.toData())
	default:
		logger.Error.Printf("Неопознанное сообщение от сервера %d", ctrl.input[6])
		ctrl.sendLive()
	}

}
