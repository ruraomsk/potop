package transport

import "sync"

var mapHistory map[int][]History
var mutex sync.Mutex

func GetHistory(level int) []History {
	mutex.Lock()
	defer mutex.Unlock()
	hi, is := mapHistory[level]
	if !is {
		return make([]History, 0)
	}
	return hi
}

func history() {
	mapHistory = make(map[int][]History)
	for {
		h := <-toHistory
		mutex.Lock()
		hi, is := mapHistory[h.Type]
		if !is {
			hi = make([]History, 0)
		}
		if len(hi) > 50 {
			hi = hi[1:]
		}
		hi = append(hi, h)
		mapHistory[h.Type] = hi
		mutex.Unlock()
	}
}
