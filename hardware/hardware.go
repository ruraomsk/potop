package hardware

import (
	"fmt"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

var HoldsCmd chan WriteHolds
var CoilsCmd chan WriteCoils
var SetWork chan int //команды управления 1 - перейти в режим управления Utopia 0- включить локальный план управления
var StateHardware = StateHard{Connect: false, Utopia: true, LastOperation: time.Unix(0, 0), Status: make([]byte, 4)}
var client *modbus.ModbusClient
var err error
var mutex sync.Mutex

func Start() {
	StateHardware.setConnect(false)
	count := 0
	HoldsCmd = make(chan WriteHolds)
	CoilsCmd = make(chan WriteCoils)
	SetWork = make(chan int)
	tickerConnect := time.NewTicker(5 * time.Second)
	tickerStatus := time.NewTicker(time.Second)
	for {
		select {
		case <-tickerConnect.C:
			if !StateHardware.GetConnect() {
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
					if (count % 100) == 0 {
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
				StateHardware.setConnect(true)
			}
		case cmd := <-SetWork:
			if cmd == 0 {
				StateHardware.setUtopia(false)
			}
			if cmd == 1 {
				StateHardware.setUtopia(true)
			}
		case <-tickerStatus.C:
			if StateHardware.GetConnect() {
				err = readStatus(StateHardware.getUtopia())
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					StateHardware.setConnect(false)
				}
			}

		case wc := <-CoilsCmd:
			logger.Debug.Printf("coils cmd %v", wc)
			if StateHardware.GetConnect() {
				err = client.WriteCoils(wc.Start, wc.Data)
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					StateHardware.setConnect(false)
				} else {
					StateHardware.setLastOperation()
				}
			}
		case wh := <-HoldsCmd:
			logger.Debug.Printf("holds cmd %v", wh)
			if StateHardware.GetConnect() {
				err = client.WriteRegisters(wh.Start, wh.Data)
				if err != nil {
					logger.Error.Print(err.Error())
					client.Close()
					StateHardware.setConnect(false)
				}
			} else {
				StateHardware.setLastOperation()
			}
		}
	}
}
func readStatus(utopia bool) error {
	mutex.Lock()
	defer mutex.Unlock()
	if !utopia {
		//utopia отключена
		StateHardware.WatchDog = 0
		err := client.WriteRegister(178, StateHardware.WatchDog)
		if err != nil {
			return fmt.Errorf("write holds 178 %s", err.Error())
		}
	}
	//Обновляем wtchdog если нужно
	if StateHardware.WatchDog > 0 {
		StateHardware.WatchDog--
		err := client.WriteRegister(178, StateHardware.WatchDog)
		if err != nil {
			return fmt.Errorf("write holds 178 %s", err.Error())
		}
	}
	//Считываем состояние направлений
	data, err := client.ReadRegisters(190, 32, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 190 32 %s", err.Error())
	}
	for i, v := range data {
		StateHardware.StatusDirs[i] = uint8(v)
	}
	//Обновляем статус КДМ в его кодах
	status, err := client.ReadRegisters(0, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 0 4 %s", err.Error())
	}
	for i, v := range status {
		StateHardware.Status[i] = uint8(v)
	}
	//Обновляем информацию о спец режимах
	coils, err := client.ReadCoils(0, 3)
	if err != nil {
		return fmt.Errorf("read coils 0 3 %s", err.Error())
	}

	StateHardware.Dark = coils[0]
	StateHardware.AllRed = coils[1]
	StateHardware.Flashing = coils[2]
	utopiacmd, err := client.ReadRegisters(175, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 175 4 %s", err.Error())
	}
	StateHardware.Tmin = int(utopiacmd[0])
	StateHardware.RealWatchDog = utopiacmd[3]
	StateHardware.MaskCommand = uint32(utopiacmd[1])<<16 | uint32(utopiacmd[2])
	return nil
}
