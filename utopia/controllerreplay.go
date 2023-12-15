package utopia

import (
	"fmt"
	"time"

	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/stat"
)

// Controller Message 190 - Status and detections: message
type StatusAndDetections struct {
	lastop    time.Time
	plan      int  //Номер плана
	TLCstatus byte //Состояние контроллера
	TLCbasic  byte //Состояние диагностики
	//The following correspondences are proposed as an example:
	//  bit00 (0x01): “watch-dog expired”
	//  bit01 (0x02): “communication error”
	//  bit02 (0x04): “conflicting signal group command”
	//  bit03 (0x08): “centralisation inhibited”
	//  bit04 (0x10): “inter-green violation”
	//  bit05 (0x20): “lamp fault”
	//  bit06 (0x40): [not used]
	//  bit07 (0x80): “extended diagnostics update”
	sensors [120]sensor
}

type sensor struct {
	counts    int //Проехало за последний период
	occupancy int
}

func (s *StatusAndDetections) ToString() string {
	var st = "undef"
	switch s.TLCstatus {
	case 0:
		st = "NO COMMAND"
	case 1:
		st = "LOCAL"
	case 2:
		st = "CENTRAL"
	case 3:
		st = "FLASHING"
	case 4:
		st = "ALL RED"
	case 5:
		st = "MANUAL"
	case 6:
		st = "DARK"
	case 7:
		st = fmt.Sprintf("RUN PLAN %d", s.plan)
	}
	var bs = ""
	if s.TLCbasic&1 != 0 {
		bs += "watch-dog expired"
	}
	if s.TLCbasic&2 != 0 {
		bs += "communication error"
	}
	if s.TLCbasic&0x4 != 0 {
		bs += "conflicting signal group command"
	}
	if s.TLCbasic&0x8 != 0 {
		bs += "centralisation inhibited"
	}
	if s.TLCbasic&0x10 != 0 {
		bs += "inter-green violation"
	}
	if s.TLCbasic&0x20 != 0 {
		bs += "lamp fault"
	}
	if s.TLCbasic&0x40 != 0 {
		bs += "not used"
	}
	if s.TLCbasic&0x80 != 0 {
		bs += "extended diagnostics update"
	}
	return fmt.Sprintf("Message 190 %s %s %s \n %v", toString(s.lastop), st, bs, s.sensors[:32])
}

func (s *StatusAndDetections) fill() {
	s.TLCbasic = hardware.GetDiagnosticUtopia()
	s.TLCstatus = hardware.GetStatusUtopia()
	s.plan = hardware.GetPlan()
	if s.plan < 1 || s.plan > 24 {
		s.plan = 0
	}
	for i := 0; i < len(s.sensors); i++ {
		s.sensors[i].counts = 0xff
		s.sensors[i].occupancy = 0xff
	}
	for i, v := range stat.GetCountValues() {
		s.sensors[i].counts = v
		if v == 0 {
			s.sensors[i].occupancy = 0
		} else {
			s.sensors[i].occupancy = 1
		}
	}
	stat.ClearCountValues()
	// logger.Debug.Printf("fill StatusAndDetections %v", s)
}

func (s *StatusAndDetections) toData() []byte {
	s.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 190, byte(s.plan), s.TLCstatus, s.TLCbasic)
	for _, v := range s.sensors {
		result = append(result, byte(v.counts), byte(v.occupancy))
	}
	return result
}
func (s *StatusAndDetections) fromData(data []byte) error {
	if data[0] != 190 {
		return fmt.Errorf("это не сообщение 190 %d", data[0])
	}
	s.lastop = time.Now()
	s.plan = int(data[1])
	s.TLCstatus = data[2]
	s.TLCbasic = data[3]
	j := 4
	for i := 0; i < 120; i++ {
		s.sensors[i].counts = int(data[j])
		j++
		s.sensors[i].occupancy = int(data[j])
		j++
	}
	return nil
}

// Controller Message 4 – Signal Group feedback, ответ от сигнальных групп
type SignalGroupFeedback struct {
	lastop time.Time
	states [64]int //Состояние светофорной группы Может принимать значения
	// Красный, красный (пешеходный) 					0х0
	// Желтый, зеленый мигающий и желтый				0х1
	// Красный, переходящий в желтый (перед зеленым)	0х2
	// Мигающий зеленый (пешеходный)					0х3
	// Зеленый, зеленый пешеходный						0х8
	// Мигающий желтый									0х9
	// Мигающий зеленый (перед желтым)					0хA
	// Ждущий зеленый 									0хB
	// Неработающий 									0хE
	// Группа не определена								0хF
}

func (s *SignalGroupFeedback) ToString() string {
	return fmt.Sprintf("Message 4 %s  % 02X", toString(s.lastop), s.states)
}

func (s *SignalGroupFeedback) fill() {
	for i := 0; i < len(s.states); i++ {
		s.states[i] = 255
	}
	for i, v := range hardware.GetStatusDirs() {
		s.states[i] = int(v)
	}
	// logger.Debug.Printf("fill SignalGroupFeedback %v", s)
}

func (s *SignalGroupFeedback) toData() []byte {
	s.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 4)
	b := make([]byte, 32)
	j := 0
	l := 4
	for _, v := range s.states {
		b[j] |= ((byte(v) & 0xf) << l)
		l -= 4
		if l < 0 {
			j++
			l = 4
		}
	}
	result = append(result, b...)
	return result
}
func (s *SignalGroupFeedback) fromData(data []byte) error {
	if data[0] != 4 {
		return fmt.Errorf("это не сообщение 4")
	}
	s.lastop = time.Now()
	j := 1
	l := 4
	for i := 0; i < len(s.states); i++ {
		s.states[i] = int((data[j] >> l) & 0xf)
		l -= 4
		if l < 0 {
			j++
			l = 4
		}
	}
	return nil
}

// Message 5 – Extended diagnostic
type ExtendedDiagnostic struct {
	lastop   time.Time
	Extrrors []ExtError
}
type ExtError struct {
	code [3]byte
}

func (e *ExtendedDiagnostic) fill() {
	e.lastop = time.Now()
	e.Extrrors = make([]ExtError, 0)
	e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{0, 0, 0}})
	e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{0, 3, 1}})
	hs := hardware.GetStateHard()
	if hs.AllRed {
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{12, 0, 0}})
	}
	switch hs.Status[0] {
	case 0:
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{0, 0, 0}})
	case 1:
		// return fmt.Sprintf("Лампа сгорела, контроль красных плата %d ключ %d",
		// 	StateHardware.Status[1], StateHardware.Status[2])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{6, 3, byte((hs.Status[1] << 4) | hs.Status[2])}})
	case 2:
		// return fmt.Sprintf("Лампа сгорела, контроль зеленых плата %d ключ %d",
		// 	StateHardware.Status[1], StateHardware.Status[2])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{6, 1, byte((hs.Status[1] << 4) | hs.Status[2])}})

	case 3:
		// return "Нет ответа от микросхемы аппаратных часов"
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{15, 0, 0}})
	case 4:
		// return "Нет сигнала от GPS приемника"
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{15, 0, 0}})
	case 5:
		// return fmt.Sprintf("Нет ответа от платы силовых ключей плата %d",
		// 	StateHardware.Status[1])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{11, 1, hs.Status[1]}})

	case 6:
		// return fmt.Sprintf("Нет ответа от платы ввода-вывода плата %d",
		// 	StateHardware.Status[1])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{11, 1, hs.Status[1]}})
	case 7:
		// return fmt.Sprintf("КЗ цепи кнопки КВП %d",
		// 	StateHardware.Status[1])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{11, 1, hs.Status[1]}})
	case 8:
		// return "версия файла конфигурации в ПЗУ не соответствует требуемой"
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{14, 1, 1}})
	case 9:
		// return "контрольная сумма файла конфигурации в ПЗУ показывает ошибку"
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{14, 1, 2}})
	case 10:
		// return fmt.Sprintf("обнаружен конфликт направлений %d",
		// 	StateHardware.Status[1])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{2, hs.Status[1], hs.Status[2]}})

	case 11:
		// return fmt.Sprintf("команда от сети обнаружен конфликт направлений %d",
		// 	StateHardware.Status[1])
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{2, hs.Status[1], hs.Status[2]}})
	case 12:
		// return "не вхождение в координацию"
		e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{15, 0, 0}})
	}
}

func (e *ExtendedDiagnostic) toData() []byte {
	e.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 5)
	for _, v := range e.Extrrors {
		for _, w := range v.code {
			result = append(result, w)
		}
	}
	return result
}
func (e *ExtendedDiagnostic) fromData(data []byte) error {
	if data[0] != 5 {
		return fmt.Errorf("это не сообщение 5")
	}
	e.lastop = time.Now()
	e.Extrrors = make([]ExtError, 0)
	if (len(data)-1)%3 != 0 {
		return fmt.Errorf("неверная структура сообщения %d", len(data))
	}
	for i := 0; i < len(data)-1; i += 3 {
		var code [3]byte
		for j := 0; j < 3; j++ {
			code[j] = data[1+i+j]
		}
	}
	return nil
}

// Message 24 – Classified counts by vehicle length
type ClassifiedCounts struct {
	lastop   time.Time
	LCLASSES int //No. vehicle length classes described in the message
	LSENSORS int //No. classification sensors described in the message
	Sensors  []Sensor
}
type Sensor struct {
	counts []int
}

func (c *ClassifiedCounts) fill() {
	classes := 5
	sensors := 10
	c.lastop = time.Now()
	c.LCLASSES = classes
	c.LSENSORS = sensors
	var sensor Sensor
	c.Sensors = make([]Sensor, 0)
	sensor.counts = make([]int, classes)
	for i := 0; i < sensors; i++ {
		c.Sensors = append(c.Sensors, sensor)
	}
}

func (c *ClassifiedCounts) toData() []byte {
	c.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 24, byte(c.LCLASSES), byte(c.LSENSORS))
	for _, v := range c.Sensors {
		for _, w := range v.counts {
			result = append(result, byte(w))
		}
	}
	return result
}
func (c *ClassifiedCounts) fromData(data []byte) error {
	if data[0] != 24 {
		return fmt.Errorf("это не сообщение 24")
	}
	c.lastop = time.Now()
	c.LCLASSES = int(data[1])
	c.LSENSORS = int(data[2])
	c.Sensors = make([]Sensor, 0)
	if (len(data) - 3) != c.LCLASSES*c.LSENSORS {
		return fmt.Errorf("неверная структура сообщения %v", data)
	}
	c.Sensors = make([]Sensor, 0)
	for i := 0; i < c.LSENSORS; i++ {
		var sensor Sensor
		sensor.counts = make([]int, 0)
		for j := 0; j < c.LCLASSES; j++ {
			sensor.counts = append(sensor.counts, int(data[3+i*c.LCLASSES+j]))
		}
	}
	return nil
}

// Message 25 – Classified counts by vehicle speed
type ClassifiedSpeeds struct {
	lastop   time.Time
	LCLASSES int //No. vehicle length classes described in the message
	LSENSORS int //No. classification sensors described in the message
	Sensors  []Sensor
}

func (c *ClassifiedSpeeds) fill() {
	classes := 5
	sensors := 5
	c.lastop = time.Now()
	c.LCLASSES = classes
	c.LSENSORS = sensors
	var sensor Sensor
	c.Sensors = make([]Sensor, 0)
	sensor.counts = make([]int, classes)
	for i := 0; i < sensors; i++ {
		c.Sensors = append(c.Sensors, sensor)
	}
}

func (c *ClassifiedSpeeds) toData() []byte {
	c.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 25, byte(c.LCLASSES), byte(c.LSENSORS))
	for _, v := range c.Sensors {
		for _, w := range v.counts {
			result = append(result, byte(w))
		}
	}
	return result
}
func (c *ClassifiedSpeeds) fromData(data []byte) error {
	if data[0] != 25 {
		return fmt.Errorf("это не сообщение 24")
	}
	c.lastop = time.Now()
	c.LCLASSES = int(data[1])
	c.LSENSORS = int(data[2])
	c.Sensors = make([]Sensor, 0)
	if (len(data) - 3) != c.LCLASSES*c.LSENSORS {
		return fmt.Errorf("неверная структура сообщения %d", len(data))
	}
	c.Sensors = make([]Sensor, 0)
	for i := 0; i < c.LSENSORS; i++ {
		var sensor Sensor
		sensor.counts = make([]int, 0)
		for j := 0; j < c.LCLASSES; j++ {
			sensor.counts = append(sensor.counts, int(data[3+i*c.LCLASSES+j]))
		}
	}
	return nil
}

// Message 1 – Bus detection
type BusDetection struct {
	lastop          time.Time
	PTid            int //PT vehicle ID
	Direction       int
	Delay           int
	PTSensorAddress int
	PTDistance      int
	PTServiceCode   int
	PTServiceShift  int
	Doors           int
}

func (b *BusDetection) fill() {
	// logger.Debug.Printf("fill BusDetection %v", b)
}

func (b *BusDetection) toData() []byte {
	b.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 1, byte(b.PTid&0xff), byte(b.PTid>>8&0xff))
	result = append(result, byte(b.Direction), byte(b.Delay), byte(b.PTSensorAddress))
	result = append(result, byte(b.PTDistance&0xff), byte(b.PTDistance>>8&0xff))
	result = append(result, byte(b.PTServiceCode&0xff), byte(b.PTServiceCode>>8&0xff))
	result = append(result, byte(b.PTServiceShift&0xff), byte(b.PTServiceShift>>8&0xff))
	result = append(result, byte(b.Doors))
	return result
}
func (b *BusDetection) fromData(data []byte) error {
	if data[0] != 1 {
		return fmt.Errorf("это не сообщение 1")
	}
	b.lastop = time.Now()
	b.PTid = int(data[1]) | int(data[2])<<8
	b.Direction = int(data[3])
	b.Delay = int(data[4])
	b.PTSensorAddress = int(data[5])
	b.PTDistance = int(data[6]) | int(data[7])<<8
	b.PTServiceCode = int(data[8]) | int(data[9])<<8
	b.PTServiceShift = int(data[10]) | int(data[11])<<8
	b.Doors = int(data[12])

	return nil
}

// Message 7 – Reply to a special command
type ReplaySpecial struct {
	lastop        time.Time
	CommandReplay int
}

func (s *ReplaySpecial) fill() {
	// logger.Debug.Printf("fill ReplaySpecial %v", s)
}

func (s *ReplaySpecial) toData() []byte {
	s.lastop = time.Now()
	result := make([]byte, 0)
	result = append(result, 7, byte(s.CommandReplay))

	return result
}
func (s *ReplaySpecial) fromData(data []byte) error {
	if data[0] != 7 {
		return fmt.Errorf("это не сообщение 7")
	}
	s.lastop = time.Now()
	s.CommandReplay = int(data[1])
	return nil
}
