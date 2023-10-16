package stat

import "time"

type OneTick struct {
	Nomber int //Номер канала
	Type   int //Тип 0-счетчик 1-скорость
	Value  Value
}
type Value struct {
	Time   time.Time //Время получения информации
	Status int       //Качество 0 - хорошее 1 обрыв
	Value  [10]int   //Значения по типу
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
	CountValues []Value
	SpeedValues []Value
	LastCount   Value
	LastSpeed   Value
}
type Chanels struct {
}
