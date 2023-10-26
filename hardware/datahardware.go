package hardware

import "sync"

// type setDirs struct {
// 	MaskDirs [2]uint16 //Маска управления вызвать направления направления
// 	WatchDog uint8     //Тайм аут управления
// }

type WriteHolds struct {
	Start uint16
	Data  []uint16
}
type WriteCoils struct {
	Start uint16
	Data  []bool
}

type StateHard struct {
	mutex      sync.Mutex
	Connect    bool      //true если есть связь с КДМ
	Dark       bool      //true если Режим ОС
	AllRed     bool      //true если Режим Кругом Красный
	Flashing   bool      //true если Режим Желтый Мигающий
	WatchDog   uint16    //Текущий Тайм аут управления
	Plan       int       //Номер исполняемого плана контроллером КДМ
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
}

func (s *StateHard) getConnect() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Connect
}
func (s *StateHard) setConnect(set bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Connect = set
}
func SetTLC(watchdog int, sgc [32]bool) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	if !state.Connect {
		return
	}
	state.WatchDog = uint16(watchdog)
	wh := WriteHolds{Start: 175, Data: make([]uint16, 4)}
	b := make([]byte, 4)
	j := 0
	l := 7
	for i := 0; i < len(sgc); i++ {
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

	if state.Dark || state.Flashing || state.AllRed {
		CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
	}

	HoldsCmd <- wh
}
func CommandUtopia(cmd int, plan int) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	if !state.Connect {
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
		if !state.Flashing {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, true}}
		}
	case 4:
		//Переход в  режим КК
		if !state.AllRed {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, true, false}}
		}
	case 6:
		//Переход в  режим ОС
		if !state.Dark {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{true, false, false}}
		}
	case 7:
		//Хочет включить план координации
		if state.Dark || state.Flashing || state.AllRed {
			CoilsCmd <- WriteCoils{Start: 0, Data: []bool{false, false, false}}
		}
		if state.WatchDog != 0 {
			state.WatchDog = 0
			HoldsCmd <- WriteHolds{Start: 178, Data: []uint16{0}}
		}
		HoldsCmd <- WriteHolds{Start: 180, Data: []uint16{uint16(plan)}}
	}
}
