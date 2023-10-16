package radar

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

var client *modbus.ModbusClient

func pusherMaster() {
	var errClient error
	var regs []uint16
	for {
		client, errClient = modbus.NewClient(&modbus.ClientConfiguration{
			URL:     fmt.Sprintf("tcp://127.0.0.1:%d", setup.Set.ModbusRadar.Port),
			Timeout: time.Second,
			Logger:  logger.Info,
		})
		if errClient != nil {
			logger.Error.Printf("Не могу создать клиента %v", errClient)
			return
		}
		client.SetUnitId(uint8(setup.Set.ModbusRadar.ID))
		for {
			errClient = client.Open()
			if errClient != nil {
				logger.Error.Printf("Не могу открыть клиента %v", errClient)
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}
		for {
			time.Sleep(1 * time.Second)
			regs, errClient = client.ReadRegisters(0, 4, modbus.HOLDING_REGISTER)
			if errClient != nil {
				logger.Error.Printf("Не смог прочитать %v ", errClient)
				break
			}
			j := 0
			k := 0
			for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
				regs[k] |= uint16(((rand.Intn(3) >> j) & 0xf))
				j += 4
				if j > 12 {
					j = 0
					k++
				}
			}
			errClient = client.WriteRegisters(0, regs)
			if errClient != nil {
				logger.Error.Printf("Не смог отправить %v %v", regs, errClient)
				break
			}
		}
		client.Close()
		time.Sleep(time.Second)
	}
}
func pusherSlave() {

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
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
	}
}
