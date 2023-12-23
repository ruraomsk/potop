package transport

import (
	"fmt"
	"time"
)

var Commander chan Command
var Requester chan Request
var Sender chan ReplayToServer
var CommandFromWeb chan Command
var RequestFromWeb chan int
var toHistory chan History

const (
	CodeCallPhase = iota
	CodeCallPlan
	CodeCallFlash
	CodeCallAllRed
	CodeCallDark
)
const (
	CodeReqStatus = iota
	CodeReqMajor
	CodeReqPlan
	CodeReqSource
	CodeReqPhase
	CodeReqSignalGroups
	CodeReqAlarm
)

func init() {
	Commander = make(chan Command)
	Requester = make(chan Request)
	Sender = make(chan ReplayToServer)
	CommandFromWeb = make(chan Command)
	RequestFromWeb = make(chan int)
	toHistory = make(chan History)
}

type History struct {
	Time    time.Time
	Type    int
	Message string
}

func (h *History) toString() string {
	return fmt.Sprintf("%s:%s", h.Time, h.Message)
}

type Command struct {
	OID   string
	Code  int
	Value int
}
type Request struct {
	OID  string
	Code int
}
type ReplayToServer struct {
	Code    int
	Elemets []Element
}

func (r *ReplayToServer) Init() {
	r.Elemets = make([]Element, 0)
}

func (r *ReplayToServer) AddTime(t time.Time) {
	r.Elemets = append(r.Elemets, Element{Type: 0, Value: uint64(t.Unix())})
}
func (r *ReplayToServer) AddInt(v int) {
	r.Elemets = append(r.Elemets, Element{Type: 1, Value: uint64(v)})
}
func (r *ReplayToServer) toString() string {
	res := receiveString(r.Code)
	for _, v := range r.Elemets {
		switch v.Type {
		case 0:
			res += " " + time.Unix(int64(v.Value), 0).Format("2006-01-02 15:04:05")
		case 1:
			res += fmt.Sprintf(",%d", v.Value)
		}
	}
	return res
}

type Element struct {
	Type  int //0 - time 1 - int
	Value uint64
}

type Define struct {
	Code int
	OID  string
}

func Transport() {
	go receiverCommands()
	go receiverRequests()
	go senderReplay()
	go history()
}
