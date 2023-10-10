package utopia

import "fmt"

// Controller Message 190 - Status and detections: message
type StatusAndDetections struct {
	plan      int  //Номер плана
	TLCstatus byte //Сосотояние контроллера
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

func (s *StatusAndDetections) toData() []byte {
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
	states [32]int //Состояние светофорной группы Может принимать значения
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

func (s *SignalGroupFeedback) toData() []byte {
	result := make([]byte, 0)
	result = append(result, 4)
	b := make([]byte, 16)
	j := 0
	l := 4
	for _, v := range s.states {
		b[0] |= ((byte(v) & 0xf) << l)
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
	Extrrors []ExtError
}
type ExtError struct {
	code [3]byte
}

func (e *ExtendedDiagnostic) new() {
	e.Extrrors = make([]ExtError, 0)
	e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{0, 1, 0}})
	e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{12, 0, 0}})
	e.Extrrors = append(e.Extrrors, ExtError{code: [3]byte{4, 3, 1}})
}

func (e *ExtendedDiagnostic) toData() []byte {
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
	LCLASSES int //No. vehicle length classes described in the message
	LSENSORS int //No. classification sensors described in the message
	Sensors  []Sensor
}
type Sensor struct {
	counts []int
}

func (c *ClassifiedCounts) new(classes int, sensors int) {
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
	LCLASSES int //No. vehicle length classes described in the message
	LSENSORS int //No. classification sensors described in the message
	Sensors  []Sensor
}

func (c *ClassifiedSpeeds) new(classes int, sensors int) {
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
	PTid            int //PT vehicle ID
	Direction       int
	Delay           int
	PTSensorAddress int
	PTDistance      int
	PTServiceCode   int
	PTServiceShift  int
	Doors           int
}

func (b *BusDetection) toData() []byte {
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
	CommandReplay int
}

func (s *ReplaySpecial) toData() []byte {
	result := make([]byte, 0)
	result = append(result, 7, byte(s.CommandReplay))

	return result
}
func (s *ReplaySpecial) fromData(data []byte) error {
	if data[0] != 7 {
		return fmt.Errorf("это не сообщение 7")
	}
	s.CommandReplay = int(data[1])
	return nil
}
