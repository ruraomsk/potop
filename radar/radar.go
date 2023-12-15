package radar

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/stat"
	"github.com/simonvetter/modbus"
)

var eh *handler
var work = false
var diapazon = 0

func GetValues() string {
	if !work {
		if setup.Set.ModbusRadar.Master {
			return "Нет связи с сервером"
		}
		return "Север еще не запущен"
	}
	eh.lock.Lock()
	defer eh.lock.Unlock()
	return fmt.Sprintf("%s %v", eh.uptime.Format("15:04:05"), eh.dates)
}
func GetWork() bool {
	return work
}
func GetStatus() string {
	if setup.Set.ModbusRadar.Master {
		return fmt.Sprintf("master for %s:%d id=%d", setup.Set.ModbusRadar.Host, setup.Set.ModbusRadar.Port, setup.Set.ModbusRadar.ID)
	} else {
		return fmt.Sprintf("slave port %d id=%d", setup.Set.ModbusRadar.Port, setup.Set.ModbusRadar.ID)
	}
}
func Radar(diap int) {
	diapazon = diap
	eh = &handler{uptime: time.Unix(0, 0)}
	if setup.Set.ModbusRadar.Master {
		if setup.Set.ModbusRadar.Debug {
			go pusherSlave()
			time.Sleep(time.Second)
		}
		go modbusMaster()
	} else {
		if setup.Set.ModbusRadar.Debug {
			go pusherMaster()
			time.Sleep(time.Second)
		}
		go modbusServer()
	}
	// go pusher()
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		eh.lock.Lock()
		eh.unpack()
		var send []stat.OneTick
		if time.Since(eh.uptime).Seconds() > 5 {
			send = badStatistics()
		} else {
			send = goodStatistics()
		}
		eh.lock.Unlock()
		for _, v := range send {
			stat.InStat <- v
		}
	}

}
func badStatistics() []stat.OneTick {
	r := make([]stat.OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
		r = append(r, stat.OneTick{Number: i, Time: t, Value: 255, Type: 0, Diap: diapazon})
	}
	return r
}
func goodStatistics() []stat.OneTick {
	r := make([]stat.OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
		r = append(r, stat.OneTick{Number: i, Time: t, Value: int(eh.dates[i]), Type: 0, Diap: diapazon})
	}
	return r
}

var server *modbus.ModbusServer
var err error

func modbusServer() {
	server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL:        fmt.Sprintf("tcp://0.0.0.0:%d", setup.Set.ModbusRadar.Port),
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, eh)
	if err != nil {
		logger.Error.Printf("Не могу создать сервер %v", err)
		return
	}

	err = server.Start()
	if err != nil {
		logger.Error.Printf("Не могу запустить сервер %v", err)
		return
	}
	work = true
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
	}

}
func modbusMaster() {
	var client *modbus.ModbusClient
	var err error

	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("tcp://%s:%d", setup.Set.ModbusRadar.Host, setup.Set.ModbusRadar.Port),
		Timeout: 5 * time.Second,
	})

	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	client.SetUnitId(uint8(setup.Set.ModbusRadar.ID))
	for {
		count := 0
		for {
			err = client.Open()
			if err != nil {
				work = false
				if count%100 == 0 {
					logger.Error.Println(err.Error())
					count++
				}
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
		logger.Info.Printf("connecting....%s:%d", setup.Set.ModbusRadar.Host, setup.Set.ModbusRadar.Port)
		work = true
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			reg16, err := client.ReadRegisters(0, uint16(len(eh.reg16)), modbus.HOLDING_REGISTER)
			if err != nil {
				work = false
				logger.Error.Printf("modbus to %s:%d %s", setup.Set.ModbusRadar.Host, setup.Set.ModbusRadar.Port, err.Error())
				client.Close()
				ticker.Stop()
				time.Sleep(5 * time.Second)
				break
			}
			eh.lock.Lock()
			for i := 0; i < len(eh.reg16); i++ {
				eh.reg16[i] = reg16[i]
			}
			eh.uptime = time.Now()
			eh.lock.Unlock()
		}
	}
}
