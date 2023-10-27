package hardware

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

var HoldsCmd chan WriteHolds
var CoilsCmd chan WriteCoils

var state = StateHard{Connect: false, Status: make([]byte, 3)}
var client *modbus.ModbusClient
var err error

func Start() {
	state.setConnect(false)
	count := 0
	HoldsCmd = make(chan WriteHolds)
	CoilsCmd = make(chan WriteCoils)
	tickerConnect := time.NewTicker(5 * time.Second)
	tickerStatus := time.NewTicker(time.Second)
	for {
		select {
		case <-tickerConnect.C:
			if !state.getConnect() {
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
				//Переводим контроллер в состояние работа
				err = client.WriteRegister(6, 1)
				if err != nil {
					logger.Error.Print(err.Error())
					continue
				}
				count = 0
				state.setConnect(true)
			}
		case <-tickerStatus.C:
			if state.getConnect() {
				err = readStatus()
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					state.setConnect(false)
				}
			}
		case wc := <-CoilsCmd:
			logger.Debug.Printf("coils cmd %v", wc)
			if state.getConnect() {
				err = client.WriteCoils(wc.Start, wc.Data)
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					state.setConnect(false)
				}
			}
		case wh := <-HoldsCmd:
			logger.Debug.Printf("holds cmd %v", wh)
			if state.getConnect() {
				err = client.WriteRegisters(wh.Start, wh.Data)
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					state.setConnect(false)
				}
			}
		}
	}
}
func readStatus() error {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	//Обновляем wtchdog если нужно
	if state.WatchDog > 0 {
		state.WatchDog--
		err := client.WriteRegister(178, state.WatchDog)
		if err != nil {
			return err
		}
	}
	//Считываем состояние направлений
	data, err := client.ReadRegisters(190, 32, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	for i, v := range data {
		state.StatusDirs[i] = uint8(v)
	}
	//Обновляем статус КДМ в его кодах
	status, err := client.ReadRegisters(0, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	for i, v := range status {
		state.Status[i] = uint8(v)
	}
	//Обновляем информацию о спец режимах
	coils, err := client.ReadCoils(0, 3)
	if err != nil {
		return err
	}

	state.Dark = coils[0]
	state.AllRed = coils[1]
	state.Flashing = coils[2]

	return nil
}
