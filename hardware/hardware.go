package hardware

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/simonvetter/modbus"
)

var HoldsCmd chan WriteHolds
var CoilsCmd chan WriteCoils
var SetWork chan int //команды управления 1 - перейти в режим управления Utopia 0- включить локальный план управления
var StateHardware = StateHard{Connect: false, Utopia: true, LastOperation: time.Unix(0, 0), Status: make([]byte, 4),
	TOOBs: make([]uint16, 32)}
var client *modbus.ModbusClient
var err error
var mutex sync.Mutex
var nowCoils map[uint16][]bool
var nowHolds map[uint16][]uint16

func Start() {
	StateHardware.setConnect(false)
	count := 0
	nowCoils = make(map[uint16][]bool)
	nowHolds = make(map[uint16][]uint16)
	HoldsCmd = make(chan WriteHolds)
	CoilsCmd = make(chan WriteCoils)
	SetWork = make(chan int)
	tickerConnect := time.NewTicker(5 * time.Second)
	tickerStatus := time.NewTicker(time.Second)
	tickerDebug := time.NewTicker(time.Second)
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
					Timeout:  3 * time.Second,
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
				} else {
					if setup.Set.Modbus.Log {
						logger.Debug.Printf("write to 6 1")
					}
				}
				count = 0
				StateHardware.setConnect(true)
				nowCoils = make(map[uint16][]bool)
				nowHolds = make(map[uint16][]uint16)
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
		case <-tickerDebug.C:
			if setup.Set.Modbus.Debug {
				fillDebugData(StateHardware.getUtopia())
			} else {
				tickerDebug.Stop()
			}
		case wc := <-CoilsCmd:
			StateHardware.setLastOperation()
			if StateHardware.GetConnect() {
				if needCoils(wc) {
					// logger.Debug.Printf("coils cmd %v", wc)
					err = client.WriteCoils(wc.Start, wc.Data)
					if err != nil {
						logger.Error.Print(err.Error())
						client.Close()
						StateHardware.setConnect(false)
					} else {
						if setup.Set.Modbus.Log {
							logger.Debug.Printf("write coils addr=%d %v", wc.Start, wc.Data)
						}
					}
				}
			}
		case wh := <-HoldsCmd:
			StateHardware.setLastOperation()
			if StateHardware.GetConnect() {
				if needHolds(wh) {
					// logger.Debug.Printf("holds cmd %v", wh)
					err = client.WriteRegisters(wh.Start, wh.Data)
					if err != nil {
						logger.Error.Print(err.Error())
						client.Close()
						StateHardware.setConnect(false)
					} else {
						if setup.Set.Modbus.Log {
							logger.Debug.Printf("write holds addr=%d % 02X", wh.Start, wh.Data)
						}
					}
				}
			}
		}
	}
}
func needCoils(wc WriteCoils) bool {
	w, is := nowCoils[wc.Start]
	if !is {
		nowCoils[wc.Start] = wc.Data
		return true
	}
	if len(w) != len(wc.Data) {
		nowCoils[wc.Start] = wc.Data
		return true
	}
	for i := 0; i < len(w); i++ {
		if w[i] != wc.Data[i] {
			nowCoils[wc.Start] = wc.Data
			return true
		}
	}
	return false
}
func needHolds(wh WriteHolds) bool {
	w, is := nowHolds[wh.Start]
	if !is {
		nowHolds[wh.Start] = wh.Data
		return true
	}
	if len(w) != len(wh.Data) {
		nowHolds[wh.Start] = wh.Data
		return true
	}
	for i := 0; i < len(w); i++ {
		if w[i] != wh.Data[i] {
			nowHolds[wh.Start] = wh.Data
			return true
		}
	}
	return false
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
	utopiacmd, err := client.ReadRegisters(175, 4, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 175 4 %s", err.Error())
	}
	StateHardware.Tmin = int(utopiacmd[0])
	StateHardware.RealWatchDog = utopiacmd[3]
	StateHardware.MaskCommand = uint32(utopiacmd[1])<<16 | uint32(utopiacmd[2])

	// //Обновляем wtchdog если нужно
	// if StateHardware.RealWatchDog > 0 {
	// 	StateHardware.RealWatchDog--
	// 	err := client.WriteRegister(178, StateHardware.RealWatchDog)
	// 	if err != nil {
	// 		return fmt.Errorf("write holds 178 %s", err.Error())
	// 	}
	// }
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

	//Обновляем источник значений для ТООВ
	source, err := client.ReadRegisters(14104, 1, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 14104 1 %s", err.Error())
	}
	StateHardware.SourceTOOB = false
	if source[0] == 1 {
		StateHardware.SourceTOOB = true
	}
	toobs, err := client.ReadRegisters(222, 32, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("read holds 222 32 %s", err.Error())
	}
	copy(StateHardware.TOOBs, toobs)
	return nil
}
func fillDebugData(utopia bool) {
	if !utopia {
		//utopia отключена
		StateHardware.WatchDog = 0
	}
	//Обновляем wtchdog если нужно
	if StateHardware.WatchDog > 0 {
		StateHardware.WatchDog--
	}
	//Считываем состояние направлений
	data := make([]uint16, 0)
	for i := 0; i < 32; i++ {
		data = append(data, uint16(rand.Intn(11)))
	}
	for i, v := range data {
		StateHardware.StatusDirs[i] = uint8(v)
	}
	//Обновляем статус КДМ в его кодах
	status := []uint16{uint16(rand.Intn(12)), uint16(rand.Intn(10)), uint16(rand.Intn(10)), uint16(rand.Intn(10))}
	for i, v := range status {
		StateHardware.Status[i] = uint8(v)
	}
	//Обновляем информацию о спец режимах
	coils := make([]bool, 3)
	for i := 0; i < 3; i++ {
		if rand.Intn(10) < 5 {
			coils = append(coils, false)
		} else {
			coils = append(coils, true)
		}
	}

	StateHardware.Dark = coils[0]
	StateHardware.AllRed = coils[1]
	StateHardware.Flashing = coils[2]

	StateHardware.Tmin = rand.Intn(50)
	StateHardware.RealWatchDog = uint16(rand.Intn(50))
	StateHardware.MaskCommand = rand.Uint32()
	// StateHardware.MaskCommand = uint32(rand.Intn(32000))<<16 | uint32(rand.Intn(32000))
}
