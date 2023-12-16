package hardware

import (
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/journal"
)

type Conflicts struct {
	Conflicts map[int]LineConflict
}
type LineConflict struct {
	Number int
	Line   []bool
}
type DefineNaps struct {
	DefineNaps map[int]DefNap
}
type DefNap struct {
	Number      int
	Type        int // 0 -не назначен 1 транспортный 2 пешеход 3 стрелка
	FlashGreen  int
	YellowFlash int
	Red         int
	RedYellow   int
	Keys        []bool
	Tiob        int // 0 выключен 1 постоянный 3 вызывной
}
type DefPhase struct {
	Number int
	Tmin   int
	Naps   []bool
}
type DefPhases struct {
	DefPhases map[int]DefPhase
}
type Plans struct {
	Plans map[int]OnePlan
}
type OnePlan struct {
	Number int
	Type   int //0 локальный план 1 координация
	Tcycle int
	Shift  int
	Lines  []Line
}
type Line struct {
	// SIMPLE = 0, 		//0 простая фаза
	// MGR,			//1 МГР фаза
	// TVP1,			//2 вызывная фаза 1
	// TVP2,			//3 вызывная фаза 2
	// TVP12,			//4 вызывная фаза 1 и 2
	// SUB_TVP1,		//5 замещающая вызывная фаза 1
	// SUB_TVP2,		//6 замещающая вызывная фаза 2
	// SUB_TVP12,		//7 замещающая вызывная фаза 1 и 2
	Type  int
	Phase int //Номер фазы
	Start int
	Stop  int
}
type Year struct {
	Year map[int]Month
}
type Month struct {
	Number int
	Days   []int
}
type Weeks struct {
	Weeks map[int]Week
}
type Week struct {
	Number int
	Days   []int
}
type Days struct {
	Days map[int]DayPlan
}

type DayPlan struct {
	Number int
	Nplans []Nplan
}
type Nplan struct {
	Start int
	Stop  int
	Plan  int
}

var configMutex sync.Mutex

type Config struct {
	Ready      bool
	Conflicts  Conflicts
	DefineNaps DefineNaps
	DefPhases  DefPhases
	RPU        OnePlan
	Plans      Plans
	Year       Year
	Weeks      Weeks
	Days       Days
}

var config Config

func (c *Conflicts) init() {
	HoldsGet <- ReadHoldsReq{Start: 10198, Lenght: 62}
	rep := <-HoldsSend
	if rep.Code != nil {
		logger.Error.Print(rep.Code.Error())
		return
	}
	c.Conflicts = make(map[int]LineConflict)
	for i := 0; i < 31; i++ {
		l := LineConflict{Number: i, Line: make([]bool, 0)}
		s := uint32(rep.Data[i*2])<<16 | uint32(rep.Data[i*2+1])
		for i := 0; i < 32; i++ {
			st := false
			if s&0x1 != 0 {
				st = true
			}
			s = s >> 1
			l.Line = append(l.Line, st)
		}
		c.Conflicts[i] = l
	}
}
func (d *DefineNaps) init() {
	d.DefineNaps = make(map[int]DefNap)
	for j := 0; j < 31; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(10260 + (j * 10)), Lenght: 10}
		rep := <-HoldsSend
		l := DefNap{Number: j, Type: int(rep.Data[0]), FlashGreen: int(rep.Data[1]), YellowFlash: int(rep.Data[2]),
			Red: int(rep.Data[3]), RedYellow: int(rep.Data[4]), Tiob: int(rep.Data[7])}
		s := uint64(rep.Data[5])<<48 | uint64(rep.Data[6])<<32 | uint64(rep.Data[7])<<16 | uint64(rep.Data[8])
		for i := 0; i < 64; i++ {
			st := false
			if s&0x1 != 0 {
				st = true
			}
			s = s >> 1
			l.Keys = append(l.Keys, st)
		}
		d.DefineNaps[j] = l
	}
}
func (d *DefPhases) init() {
	d.DefPhases = make(map[int]DefPhase)
	for j := 0; j < 31; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(10580 + (j * 3)), Lenght: 3}
		rep := <-HoldsSend
		l := DefPhase{Number: j, Tmin: int(rep.Data[0]), Naps: make([]bool, 0)}
		s := uint32(rep.Data[1])<<16 | uint32(rep.Data[2])
		for i := 0; i < 32; i++ {
			st := false
			if s&0x1 != 0 {
				st = true
			}
			s = s >> 1
			l.Naps = append(l.Naps, st)
		}
		d.DefPhases[j] = l
	}
}
func (p *OnePlan) init() {
	HoldsGet <- ReadHoldsReq{Start: uint16(10676), Lenght: 99}
	rep := <-HoldsSend
	p.Number = 0
	p.Type = int(rep.Data[0])
	p.Tcycle = int(rep.Data[1])
	p.Shift = int(rep.Data[2])
	p.Lines = make([]Line, 0)
	pos := 3
	for i := 0; i < 24; i++ {
		p.Lines = append(p.Lines, Line{Phase: int(rep.Data[pos]), Type: int(rep.Data[pos+1]), Start: int(rep.Data[pos+2]), Stop: int(rep.Data[pos+3])})
		pos += 4
	}
}
func (p *Plans) init() {
	p.Plans = make(map[int]OnePlan)
	for j := 0; j < 24; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(10775 + (j * 99)), Lenght: 99}
		rep := <-HoldsSend
		l := OnePlan{Number: j, Type: int(rep.Data[0]), Tcycle: int(rep.Data[1]), Shift: int(rep.Data[2]), Lines: make([]Line, 0)}
		pos := 3
		for i := 0; i < 24; i++ {
			l.Lines = append(l.Lines, Line{Phase: int(rep.Data[pos]), Type: int(rep.Data[pos+1]), Start: int(rep.Data[pos+2]), Stop: int(rep.Data[pos+3])})
			pos += 4
		}
		p.Plans[j] = l
	}

}
func (y *Year) init() {
	y.Year = make(map[int]Month)
	for j := 0; j < 12; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(13151 + (j * 31)), Lenght: 31}
		rep := <-HoldsSend
		l := Month{Number: j, Days: make([]int, 0)}
		for i := 0; i < 31; i++ {
			l.Days = append(l.Days, int(rep.Data[i]))
		}
		y.Year[j] = l
	}
}
func (w *Weeks) init() {
	w.Weeks = make(map[int]Week)
	for j := 0; j < 12; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(13523 + (j * 7)), Lenght: 7}
		rep := <-HoldsSend
		l := Week{Number: j, Days: make([]int, 0)}
		for i := 0; i < 7; i++ {
			l.Days = append(l.Days, int(rep.Data[i]))
		}
		w.Weeks[j] = l
	}
}
func (d *Days) init() {
	d.Days = make(map[int]DayPlan)
	for j := 0; j < 12; j++ {
		HoldsGet <- ReadHoldsReq{Start: uint16(13607 + (j * 36)), Lenght: 36}
		rep := <-HoldsSend
		l := DayPlan{Number: j, Nplans: make([]Nplan, 0)}
		pos := 0
		for i := 0; i < 12; i++ {
			l.Nplans = append(l.Nplans, Nplan{Start: int(rep.Data[pos]), Stop: int(rep.Data[pos+1]), Plan: int(rep.Data[pos+2])})
			pos += 3
		}
		d.Days[j] = l
	}
}
func GetConfig() Config {
	configMutex.Lock()
	defer configMutex.Unlock()
	return config
}
func getconfig() {
	HoldsCmd <- WriteHolds{Start: 186, Data: []uint16{1}}
	configMutex.Lock()
	config.Conflicts.init()
	config.DefineNaps.init()
	config.DefPhases.init()
	config.Plans.init()
	config.Year.init()
	config.Weeks.init()
	config.Days.init()
	config.RPU.init()
	configMutex.Unlock()
	config.Ready = true
	HoldsCmd <- WriteHolds{Start: 186, Data: []uint16{0}}
}
func configure() {
	config = Config{Ready: false}
	for {
		if !StateHardware.GetConnect() {
			configMutex.Lock()
			config.Ready = false
			configMutex.Unlock()
			time.Sleep(1 * time.Second)
			continue
		}
		getconfig()
		journal.SendMessage(1, "Привязки прочитаны")
		for StateHardware.GetConnect() {
			time.Sleep(10 * time.Second)
			continue
		}
	}
}
