package stat

import (
	"fmt"
	"sync"
	"time"
)

type OneTick struct {
	Number int       //Номер канала
	Time   time.Time //Время получения информации
	Type   int       //Тип 0-счетчик 1-скорость
	Diap   int       //Диапазон
	Value  int
}
type Value struct {
	Status int     //Качество 0 - хорошее 1 обрыв
	Value  [10]int //Значения по типу
	//Если кол-во
	// Class			Description
	// I				Up to 2,5 meters
	// II				from 2,5 to 5 meters
	// III				from 5 to 8 meters
	// IV				from 8 to 11 meters
	// V				from 11 to 14 meters
	// VI				from 14 to 17 meters
	// VII				from 17 to 22 meters
	// VIII				Over 22 meters
	// IX				[not used]
	// X				Not classified
	// Если скорость
	// Class			Description
	// I				Up to 15 km/h
	// II				from 15 to 30 km/h
	// III				from 30 to 50 km/h
	// IV				from 50 to 60 km/h
	// V				from 60 to 70 km/h
	// VI				from 70 to 80 km/h
	// VII				from 80 to 100 km/h
	// VIII				from 100 to 130 km/h
	// IX				Over 130 km/h
	// X				Not classified

}
type OneChanel struct {
	Number      int
	CountValues Value
	SpeedValues Value
	LastCount   Value
	LastSpeed   Value
}
type Chanels struct {
	mutex   sync.Mutex
	counts  int //Кол-во датчиков
	diaps   int //Кол-во диапазонов
	chanels map[int]*OneChanel
}

func (v *Value) clear() {
	v.Status = 1
	for i := 0; i < len(v.Value); i++ {
		v.Value[i] = 0xff
	}
}
func (v *Value) add(t OneTick) error {
	if t.Diap < 0 || t.Diap > 9 {
		return fmt.Errorf("неверный номер диапазона %d", t.Diap)
	}
	if t.Value == 0xff {
		v.Status = 1
	} else {
		v.Status = 0
		if v.Value[t.Diap] == 0xff {
			v.Value[t.Diap] = t.Value
		} else {
			v.Value[t.Diap] += t.Value
		}
	}
	return nil
}

func (o *OneChanel) clear() {
	o.CountValues.clear()
	o.SpeedValues.clear()
	o.LastCount.clear()
	o.LastSpeed.clear()
}
func (c *Chanels) clear(chanels int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.counts = chanels
	c.diaps = 10
	for i := 0; i < chanels; i++ {
		v := OneChanel{Number: i}
		v.clear()
		c.chanels[i] = &v
	}
}
func (c *Chanels) add(t OneTick) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	o, is := c.chanels[t.Number]
	if !is {
		return fmt.Errorf("нет такого канала %d", t.Number)
	}
	return o.add(t)
}
func (o *OneChanel) add(t OneTick) error {
	switch t.Type {
	case 0: //Кол-во
		return o.CountValues.add(t)
	case 1: //Скорость
		return o.SpeedValues.add(t)
	default:
		return fmt.Errorf("ошибка типа сообщения %d", t.Type)
	}
}
