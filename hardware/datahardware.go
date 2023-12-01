package hardware

import (
	"fmt"
	"time"

	"github.com/ruraomsk/potop/setup"
)

type WriteHolds struct {
	Start uint16
	Data  []uint16
}
type WriteCoils struct {
	Start uint16
	Data  []bool
}

type StateHard struct {
	Utopia        bool //true управляение утопией false - локальное управление
	LastOperation time.Time
	Connect       bool   //true если есть связь с КДМ
	Dark          bool   //true если Режим ОС
	AllRed        bool   //true если Режим Кругом Красный
	Flashing      bool   //true если Режим Желтый Мигающий
	SourceTOOB    bool   //true если Источник времени отсчета внешний
	WatchDog      uint16 //Текущий Тайм аут управления
	Plan          int    //Номер исполняемого плана контроллером КДМ
	// typedef enum {					//Идентификаторы событий в логе аварий и в регистре событий
	// 	ALL_IS_GOOD = 0,			//Все хорошо, нет предупреждений
	// 	LOW_CURRENT_RED_LAMP,			//Ток через открытый ключ меньше минимального - лампа сгорела, применяется при контроле красных
	// 	NOT_ALLOWED_VOLTAGE_GREEN_OUT, 	//Обнаружено напряжение на закрытом ключе, применяется при контроле зеленых
	// 	NO_CLOCK,				//Нет ответа от микросхемы аппаратных часов
	// 	NO_GPS,					//Нет сигнала от GPS приемника
	// 	NO_POWER_BOARD,			//Нет ответа от платы силовых ключей
	// 	NO_IO_BOARD,				//Нет ответа от платы ввода-вывода
	// 	SHORT_CIRQUIT_KVP			//КЗ цепи кнопки КВП
	// 	WRONG_FILE_VER,		//версия файла конфигурации в ПЗУ не соответствует требуемой
	// 	WRONG_FILE_CRC			//контрольная сумма файла конфигурации в ПЗУ показывает ошибку
	// 	DIRECTIONS_CONFLICT		//обнаружен конфликт направлений
	// 	DC_DIRECTIONS_CONFLICT		//при вызове направлений по сети обнаружен конфликт направлений, вызов отклонен
	// 	NOT_ENTERING_COORDINATION		//не вхождение в координацию
	// }EventId;

	// для событий LOW_CURRENT_RED_LAMP и NOT_ALLOWED_VOLTAGE_GREEN_OUT, S1 содержит номер платы, S2 – номер ключа на плате;
	// для событий NO_POWER_BOARD и NO_IO_BOARD, S1 содержит номер платы;
	// для события SHORT_CIRQUIT_KVP, S1 содержит номер кнопки;
	// для событий DIRECTIONS_CONFLICT и DC_DIRECTIONS_CONFLICT, S1 содержит номер конфликтующего направления
	// для других событий описания не используются.
	Status     []byte    //Статус КДМ в его кодировке
	StatusDirs [32]uint8 //Статусы состояния по направлениям
	//   OFF = 0, //все сигналы выключены
	//   DEACTIV_YELLOW=1, //направление перешло в неактивное состояние, желтый после зеленого
	//   DEACTIV_RED=2, //направление перешло в неактивное состояние, красный
	//   ACTIV_RED=3, //направление перешло в активное состояние, красный
	//   ACTIV_REDYELLOW=4, //направление перешло в активное состояние, красный c желтым
	//   ACTIV_GREEN=5, //направление перешло в активное состояние, зеленый
	//   UNCHANGE_GREEN=6, //направление не меняло свое состояние, зеленый
	//   UNCHANGE_RED=7, //направление не меняло свое состояние, красный
	//   GREEN_BLINK=8, //зеленый мигающий сигнал
	//   ZM_YELLOW_BLINK=9, //желтый мигающий в режиме ЖМ
	//   OS_OFF=10,	//сигналы выключены в режиме ОС
	//   UNUSED=11 //неиспользуемое направление
	Tmin         int      //Последнее заданное Тмин вызвать направления
	MaskCommand  uint32   //Последняя маска
	RealWatchDog uint16   //Остаток watchdog
	TOOBs        []uint16 //Счетчики по направлениям
}

func (s *StateHard) GetConnect() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return s.Connect
}
func (s *StateHard) GetConnectUtopia() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return s.Utopia
}

func (s *StateHard) setConnect(set bool) {
	mutex.Lock()
	defer mutex.Unlock()
	s.Connect = set
}
func (s *StateHard) getUtopia() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return s.Utopia
}
func (s *StateHard) setUtopia(set bool) {
	mutex.Lock()
	defer mutex.Unlock()
	s.Utopia = set
}

func (s *StateHard) setLastOperation() {
	mutex.Lock()
	defer mutex.Unlock()
	s.LastOperation = time.Now()
}
func SetSignalCountDown(counts [64]byte) {
	if !StateHardware.SourceTOOB {
		setExternalSourceTOOB()
	}
	wh := WriteHolds{Start: 222, Data: make([]uint16, 32)}
	for i := 0; i < len(wh.Data); i++ {
		wh.Data[i] = uint16(counts[i])
	}
	// logger.Debug.Printf("%v ", wh)
	HoldsCmd <- wh
}
func setInternalSourceTOOB() {
	wh1 := WriteHolds{Start: 186, Data: []uint16{1}}
	HoldsCmd <- wh1
	wh2 := WriteHolds{Start: 14104, Data: []uint16{0}}
	HoldsCmd <- wh2
	wh3 := WriteHolds{Start: 186, Data: []uint16{2}}
	HoldsCmd <- wh3
}
func setExternalSourceTOOB() {
	wh1 := WriteHolds{Start: 186, Data: []uint16{1}}
	HoldsCmd <- wh1
	wh2 := WriteHolds{Start: 14104, Data: []uint16{1}}
	HoldsCmd <- wh2
	wh3 := WriteHolds{Start: 186, Data: []uint16{2}}
	HoldsCmd <- wh3
}
func SetTLC(watchdog int, sgc [64]bool) {
	if !StateHardware.Connect {
		return
	}
	StateHardware.WatchDog = uint16(watchdog)
	wh := WriteHolds{Start: 175, Data: make([]uint16, 4)}
	m := uint32(0)
	for i := 0; i < 31; i++ {
		if sgc[i] {
			m |= 0x80000000
		}
		m = m >> 1
	}
	if sgc[31] {
		m |= 0x80000000
	}

	wh.Data[1] = uint16(m >> 16 & 0xffff)
	wh.Data[2] = uint16(m & 0xffff)
	wh.Data[0] = uint16(setup.Set.Utopia.Tmin)
	if m == 0 {
		wh.Data[0] = 0
	}
	wh.Data[3] = uint16(watchdog)
	// logger.Debug.Printf("%x %v", m, sgc)
	if StateHardware.Dark || StateHardware.Flashing || StateHardware.AllRed {
		CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
	}

	HoldsCmd <- wh
}
func GetPlan() int {
	mutex.Lock()
	defer mutex.Unlock()

	return StateHardware.Plan
}
func IsConnectedKDM() bool {
	return StateHardware.GetConnect()
}
func GetStatusDirs() []uint8 {
	mutex.Lock()
	defer mutex.Unlock()
	result := make([]uint8, 0)
	var b uint8
	for _, v := range StateHardware.StatusDirs {
		switch v {
		case 0:
			//   OFF = 0, //все сигналы выключены
			b = 0xE
		case 1:
			//   DEACTIV_YELLOW=1, //направление перешло в неактивное состояние, желтый после зеленого
			b = 0x1
		case 2:
			//   DEACTIV_RED=2, //направление перешло в неактивное состояние, красный
			b = 0x0
		case 3:
			//   ACTIV_RED=3, //направление перешло в активное состояние, красный
			b = 0x0
		case 4:
			//   ACTIV_REDYELLOW=4, //направление перешло в активное состояние, красный c желтым
			b = 0x2
		case 5:
			//   ACTIV_GREEN=5, //направление перешло в активное состояние, зеленый
			b = 0x8
		case 6:
			//   UNCHANGE_GREEN=6, //направление не меняло свое состояние, зеленый
			b = 0x8
		case 7:
			//   UNCHANGE_RED=7, //направление не меняло свое состояние, красный
			b = 0x0
		case 8:
			//   GREEN_BLINK=8, //зеленый мигающий сигнал
			b = 0xA
		case 9:
			//   ZM_YELLOW_BLINK=9, //желтый мигающий в режиме ЖМ
			b = 0x9
		case 10:
			//   OS_OFF=10,	//сигналы выключены в режиме ОС
			b = 0xe
		case 11:
			//   UNUSED=11 //неиспользуемое направление
			b = 0xf
		default:
			b = 0xe
		}
		result = append(result, b)
	}
	return result
}

func GetStatusUtopia() byte {
	mutex.Lock()
	defer mutex.Unlock()
	if StateHardware.Dark {
		return 6
	}
	if StateHardware.Flashing {
		return 3
	}
	if StateHardware.AllRed {
		return 4
	}
	if StateHardware.WatchDog == 0 {
		return 1
	}
	if StateHardware.WatchDog != 0 {
		return 2
	}
	return 5
}
func GetDiagnosticUtopia() byte {
	mutex.Lock()
	defer mutex.Unlock()
	result := byte(0)
	//bit01 (0x02): “communication error”
	if !StateHardware.Connect {
		result |= 2
		return result
	}
	if StateHardware.WatchDog == 0 {
		result |= 1
	}
	if StateHardware.Status[0] == 0 {
		return result
	}
	// bit02 (0x04): “conflicting signal group command”
	if StateHardware.Status[0] == 10 || StateHardware.Status[0] == 11 {
		result |= 0x04
	}
	//  bit03 (0x08): “centralisation inhibited”
	//  bit04 (0x10): “inter-green violation”
	if StateHardware.Status[0] == 2 {
		result |= 0x10
	}
	//  bit05 (0x20): “lamp fault”
	if StateHardware.Status[0] == 1 {
		result |= 0x20
	}
	//  bit06 (0x40): [not used]
	//  bit07 (0x80): “extended diagnostics update”
	if StateHardware.Status[0] >= 3 && StateHardware.Status[0] <= 9 {
		result |= 0x80
	}
	return result
}
func CommandUtopia(cmd int, plan int) {
	// mutex.Lock()
	// defer mutex.Unlock()
	if !StateHardware.Connect {
		return
	}
	if StateHardware.SourceTOOB {
		setInternalSourceTOOB()
	}
	switch cmd {
	case 0:
		//no command
		return
	case 1:
		//Переход в локальный режим
		CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
		StateHardware.WatchDog = uint16(0)
		wh := WriteHolds{Start: 175, Data: make([]uint16, 4)}
		HoldsCmd <- wh
	case 3:
		//Переход в  режим ЖМ
		if !StateHardware.Flashing {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, true}}
		}
	case 4:
		//Переход в  режим КК
		if !StateHardware.AllRed {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, true, false}}
		}
	case 6:
		//Переход в  режим ОС
		if !StateHardware.Dark {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{true, false, false}}
		}
	case 7:
		//Хочет включить план координации
		if StateHardware.Dark || StateHardware.Flashing || StateHardware.AllRed {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
		}
		if StateHardware.WatchDog != 0 {
			StateHardware.WatchDog = 0
			HoldsCmd <- WriteHolds{Start: 178, Data: []uint16{0}}
		}
		HoldsCmd <- WriteHolds{Start: 180, Data: []uint16{uint16(plan)}}
	}
}
func GetStateHard() StateHard {
	mutex.Lock()
	defer mutex.Unlock()
	return StateHardware
}
func GetError() string {
	mutex.Lock()
	defer mutex.Unlock()
	switch StateHardware.Status[0] {
	case 0:
		return "Нет ошибок"
	case 1:
		return fmt.Sprintf("Лампа сгорела, контроль красных плата %d ключ %d",
			StateHardware.Status[1], StateHardware.Status[2])
	case 2:
		return fmt.Sprintf("Лампа сгорела, контроль зеленых плата %d ключ %d",
			StateHardware.Status[1], StateHardware.Status[2])
	case 3:
		return "Нет ответа от микросхемы аппаратных часов"
	case 4:
		return "Нет сигнала от GPS приемника"
	case 5:
		return fmt.Sprintf("Нет ответа от платы силовых ключей плата %d",
			StateHardware.Status[1])
	case 6:
		return fmt.Sprintf("Нет ответа от платы ввода-вывода плата %d",
			StateHardware.Status[1])
	case 7:
		return fmt.Sprintf("КЗ цепи кнопки КВП %d",
			StateHardware.Status[1])
	case 8:
		return "версия файла конфигурации в ПЗУ не соответствует требуемой"
	case 9:
		return "контрольная сумма файла конфигурации в ПЗУ показывает ошибку"
	case 10:
		return fmt.Sprintf("обнаружен конфликт направлений %d",
			StateHardware.Status[1])
	case 11:
		return fmt.Sprintf("команда от сети обнаружен конфликт направлений %d",
			StateHardware.Status[1])
	case 12:
		return "не вхождение в координацию"
	}

	return fmt.Sprintf("Ошибка %v", StateHardware.Status)
}
