package traffic

import (
	"encoding/xml"
	"net"
	"strconv"

	"github.com/ruraomsk/ag-server/logger"
)

func sendGetTime(soc net.Conn) error {
	m := Message{Type: "GetTime"}
	s, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return sendString(soc, string(s))
}
func sendUpTime(soc net.Conn) error {
	m := Message{Type: "GetUpTime"}
	s, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return sendString(soc, string(s))
}
func sendConnectionInformation(soc net.Conn) error {
	m := Message{Type: "GetConnectionInformation"}
	s, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return sendString(soc, string(s))
}
func sendProductInformation(soc net.Conn) error {
	m := Message{Type: "ProductInformation"}
	s, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return sendString(soc, string(s))
}
func sendSubscription(soc net.Conn) error {
	m := Subscription{Type: "Subscription"}
	m.Body.Subscription = make([]Subs, 0)
	m.Body.Subscription = append(m.Body.Subscription, Subs{Type: "CarExit", Action: "Subscribe"})
	m.Body.Subscription = append(m.Body.Subscription, Subs{Type: "CarEnter", Action: "Subscribe"})
	s, err := xml.Marshal(m)
	if err != nil {
		return err
	}
	return sendString(soc, string(s))
}
func replayConnectionInformation(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c ComnectionInformation
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	connectionInformation = c
	// logger.Debug.Printf("%v", c)
}
func replayTime(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c Time
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	getTime = c
	// logger.Debug.Printf("%v", c)
}
func replayUpTime(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c UpTime
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	upTime = c
	// logger.Debug.Printf("%v", c)
}
func replayProductInformation(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c ProductInformation
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	productInformation = c
	// logger.Debug.Printf("%v", c)
}
func replayImage(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c Image
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	image = c
	// logger.Debug.Printf("%v", c)
}
func replaySubcription(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c Subscription
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	// logger.Debug.Printf("subscription %v", c)
}
func replayEvent(buffer string) {
	mutex.Lock()
	defer mutex.Unlock()
	var c CarMoved
	err := xml.Unmarshal([]byte(buffer), &c)
	if err != nil {
		return
	}
	if c.Body.Type != "CarExit" {
		return
	}
	zoneid, _ := strconv.Atoi(c.Body.ZoneId)
	zoneid--
	// cartype, _ := strconv.Atoi(c.Body.CarType)
	if zoneid < 0 || zoneid >= len(datas) {
		logger.Debug.Printf("event %v", c)
		return
	}
	datas[zoneid]++
}
