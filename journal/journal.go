package journal

import (
	"strings"
	"sync"
	"time"
)

var maxCounts = 100
var Logger chan LogMessage
var maps map[int]Header
var mutex sync.Mutex

type LogMessage struct {
	Level   int
	Message string
}
type Point struct {
	Time  time.Time
	Value string
}

type Header struct {
	Now    string
	Points []Point
}

func add(level int, data string) {
	mutex.Lock()
	defer mutex.Unlock()
	h, is := maps[level]
	if !is {
		h = Header{Now: data, Points: make([]Point, 0)}
		h.Points = append(h.Points, Point{Time: time.Now(), Value: data})
		maps[level] = h
		return
	}
	h.add(data)
	maps[level] = h
}
func (h *Header) add(data string) {
	if strings.Compare(h.Now, data) == 0 {
		return
	}
	if len(h.Points) == maxCounts {
		h.Points = h.Points[1:]
	}
	h.Now = data
	h.Points = append(h.Points, Point{Time: time.Now(), Value: data})
}
func init() {
	Logger = make(chan LogMessage, 100)
	maps = make(map[int]Header)
}
func LoggerStart() {
	for {
		m := <-Logger
		add(m.Level, m.Message)
	}
}
func GetMessages(level int) []Point {
	mutex.Lock()
	defer mutex.Unlock()
	h, is := maps[level]
	if !is {
		return make([]Point, 0)
	}
	return h.Points
}
func SendMessage(level int, data string) {
	Logger <- LogMessage{Level: level, Message: data}
}
