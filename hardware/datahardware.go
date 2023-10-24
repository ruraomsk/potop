package hardware

import "sync"

type setDirs struct {
	MaskDirs uint32 //Маска управления вызвать направления направления
	WatchDog uint8  //Тайм аут управления
}

type StateHard struct {
	mutex      sync.Mutex
	Connect    bool      //true если есть связь с КДМ
	Dark       bool      //true если Режим ОС
	AllRed     bool      //true если Режим Кругом Красный
	Flashing   bool      //true если Режим Желтый Мигающий
	Number     int       //Серийный номер КДМ
	MaskDirs   uint32    //Маска управления вызвать направления направления
	WatchDog   uint8     //Тайм аут управления
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
	TLCChan   chan setDirs
	CmdUtopia chan int //Команда от утопии 1 - перйти в локальный режим 2 - команда из цента 3 - ЖМ 6 - ОС 7- запустить план координации
	Plan      chan int
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
	set := setDirs{WatchDog: uint8(watchdog)}
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
	set.MaskDirs = uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
	state.TLCChan <- set
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
		state.CmdUtopia <- 1
	case 3:
		//Переход в  режим ЖМ
		state.CmdUtopia <- 3
	case 6:
		//Переход в  режим ОС
		state.CmdUtopia <- 6
	case 7:
		state.Plan <- plan
	}
}
