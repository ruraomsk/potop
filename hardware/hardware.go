package hardware

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

var state StateHard
var client *modbus.ModbusClient
var err error

func Start() {
	state.setConnect(false)
	count := 0
	ticker := time.NewTicker(time.Second)
	for !state.getConnect() {
		select {
		case <-ticker.C:
			// for an RTU (serial) device/bus
			client, err = modbus.NewClient(&modbus.ClientConfiguration{
				URL:      setup.Set.Modbus.Device,         //"rtu:///dev/ttyUSB0",
				Speed:    uint(setup.Set.Modbus.BaudRate), //19200,                   // default
				DataBits: 8,                               // default, optional
				Parity:   modbus.PARITY_NONE,              // default, optional
				StopBits: 2,                               // default if no parity, optional
				Timeout:  300 * time.Millisecond,
			})

			if err != nil {
				logger.Error.Printf("modbus %v", err.Error())
				continue
			}
			client.SetUnitId(uint8(setup.Set.Modbus.UId))
			err = client.Open()
			if err != nil {
				if (count % 1000) == 0 {
					logger.Error.Printf("modbus open %v", err.Error())
				}
				count++
				continue
			}
			count = 0
			state.setConnect(true)
			ticker.Stop()
		}
		workModbus()
		logger.Error.Printf("Завершили обмен с ModBus")
		state.setConnect(false)
	}
}
func workModbus() {

}
