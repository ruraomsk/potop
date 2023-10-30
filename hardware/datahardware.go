package hardware

import (
	"time"
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
	Tmin         int    //Последнее заданное Тмин вызвать направления
	MaskCommand  uint32 //Последняя маска
	RealWatchDog uint16 //Остаток watchdog
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

func SetTLC(watchdog int, sgc [64]bool) {
	mutex.Lock()
	defer mutex.Unlock()
	if !StateHardware.Connect {
		return
	}
	StateHardware.WatchDog = uint16(watchdog)
	wh := WriteHolds{Start: 175, Data: make([]uint16, 4)}
	b := make([]byte, 4)
	j := 0
	l := 7
	for i := 0; i < 32; i++ {
		d := 0
		if sgc[i] {
			d = 1
		}
		d = d << l
		b[j] |= byte(d)
		l--
		if l < 0 {
			j++
			l = 7
		}
	}

	wh.Data[1] = uint16(b[3])<<8 | uint16(b[2])
	wh.Data[2] = uint16(b[1])<<8 | uint16(b[0])
	wh.Data[0] = uint16(watchdog)
	wh.Data[3] = uint16(watchdog)

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
func GetStatusDirs() []uint8 {
	mutex.Lock()
	defer mutex.Unlock()
	result := make([]uint8, 0)
	var b uint8
	for _, v := range StateHardware.StatusDirs {
		switch v {
		//   OFF = 0, //все сигналы выключены
		case 0:
			b = 0xE
		//   DEACTIV_YELLOW=1, //направление перешло в неактивное состояние, желтый после зеленого
		case 1:
			b = 0x1
			//   DEACTIV_RED=2, //направление перешло в неактивное состояние, красный
		case 2:
			b = 0x0
			//   ACTIV_RED=3, //направление перешло в активное состояние, красный
		case 3:
			b = 0x2
			//   ACTIV_REDYELLOW=4, //направление перешло в активное состояние, красный c желтым
		case 4:
			b = 0x2
			//   ACTIV_GREEN=5, //направление перешло в активное состояние, зеленый
		case 5:
			b = 0x11
			//   UNCHANGE_GREEN=6, //направление не меняло свое состояние, зеленый
		case 6:
			b = 0x11
			//   UNCHANGE_RED=7, //направление не меняло свое состояние, красный
		case 7:
			b = 0x0
			//   GREEN_BLINK=8, //зеленый мигающий сигнал
		case 8:
			b = 0xA
			//   ZM_YELLOW_BLINK=9, //желтый мигающий в режиме ЖМ
		case 9:
			b = 0x3
			//   OS_OFF=10,	//сигналы выключены в режиме ОС
		case 10:
			b = 0xE
			//   UNUSED=11 //неиспользуемое направление
		case 11:
			b = 0xf
		default:
			b = 0xF
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
	mutex.Lock()
	defer mutex.Unlock()
	if !StateHardware.Connect {
		return
	}
	switch cmd {
	case 0:
		//no command
		return
	case 1:
		//Переход в локальный режим
		CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
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
