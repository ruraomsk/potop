package diagramm

import (
	"sync"
	"time"

	"github.com/ruraomsk/potop/hardware"
)

const maxCount = 100

var mutex sync.Mutex

type Pdiag struct {
	Time time.Time
	Data []int
}
type ViewDiag struct {
	View   bool
	Chanel int
	Data   []string
}

var diags []Pdiag

func DiagrammStart() {
	time.Sleep(time.Second)
	diags = make([]Pdiag, 0)
	for {
		time.Sleep(time.Second)
		hs := hardware.GetStateHard()
		if !hs.Connect {
			diags = make([]Pdiag, 0)
			continue
		}
		m := make([]int, 0)
		for _, v := range hs.StatusDirs {
			m = append(m, int(v))
		}
		if len(diags) == maxCount {
			diags = diags[1:]
		}
		mutex.Lock()
		diags = append(diags, Pdiag{Time: time.Now(), Data: m})
		mutex.Unlock()
	}
}

func GetDiagramm(start time.Time, finish time.Time) ([]ViewDiag, []time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	res := make([]ViewDiag, 0)
	times := make([]time.Time, 0)
	for i := 0; i < 32; i++ {
		res = append(res, ViewDiag{Chanel: i, Data: make([]string, 0), View: true})
	}
	for _, v := range diags {
		if v.Time.After(start) && v.Time.Before(finish) {
			times = append(times, v.Time)
			for i, k := range v.Data {
				view := true
				ds := " "

				switch k {
				case 0:
					ds = "ВВ"
				case 1:
					ds = "ЖЖ"
				case 2:
					ds = "КК"
				case 3:
					ds = "КК"
				case 4:
					ds = "КЖ"
				case 5:
					ds = "ЗЗ"
				case 6:
					ds = "ЗЗ"
				case 7:
					ds = "КК"
				case 8:
					ds = "Зз"
				case 9:
					ds = "ЖМ"
				case 10:
					ds = "ОС"
				case 11:
					ds = "__"
					view = false
				default:
					view = false
				}
				res[i].Data = append(res[i].Data, ds)
				if !view {
					res[i].View = view
				}
			}
		}
	}
	return res, times
}
