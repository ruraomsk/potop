package setup

var (
	Set    *Setup
	ExtSet *ExtSetup
)

type Setup struct {
	LogPath string `toml:"logpath"`
	Id      int    `toml:"id"`
	Modbus  Modbus `toml:"modbus" json:"modbus"`
	Utopia  Utopia `toml:"utopia" json:"utopia"`
	STCIP   STCIP  `toml:"stcip" json:"stcip"`

	TrafficData TrafficData `toml:"trafficdata" json:"trafficdata"`
	ModbusRadar ModbusRadar `toml:"modbusradar" json:"modbusradar"`
}
type ExtSetup struct {
	Modbus      Modbus      `toml:"modbus" json:"modbus"`
	Utopia      Utopia      `toml:"utopia" json:"utopia"`
	STCIP       STCIP       `toml:"stcip" json:"stcip"`
	TrafficData TrafficData `toml:"trafficdata" json:"trafficdata"`
	ModbusRadar ModbusRadar `toml:"modbusradar" json:"modbusradar"`
}
type STCIP struct {
	Run    bool
	Debug  bool   `toml:"debug" json:"debug"`
	Host   string `toml:"host" json:"host"`
	Port   int    `toml:"port" json:"port"`
	Listen int    `toml:"listen" json:"listen"`
}

type ModbusRadar struct {
	Work    bool   `toml:"work" json:"work"`
	Master  bool   `toml:"master" json:"master"`
	Debug   bool   `toml:"debug" json:"debug"`
	Host    string `toml:"host" json:"host"`
	Port    int    `toml:"port" json:"port"`
	ID      int    `toml:"id" json:"id"`
	Chanels int    `toml:"chanels" json:"chanels"`
	Diaps   int    `toml:"diaps" json:"diaps"`
	Diap    int    `toml:"diap" json:"diap"`
}

type Utopia struct {
	Run         bool
	Device      string `toml:"device" json:"device"`
	BaudRate    int    `toml:"baudrate" json:"baudrate"`
	Parity      string `toml:"parity" json:"parity"`
	UId         int    `toml:"uid" json:"uid"`
	Debug       bool   `toml:"debug" json:"debug"`
	LostControl int    `toml:"lostControl" json:"lostControl"`
	Recode      bool   `toml:"recode" json:"recode"`
	Replay      bool   `toml:"replay" json:"replay"`
}
type TrafficData struct {
	Work    bool   `toml:"work" json:"work"`
	Debug   bool   `toml:"debug" json:"debug"`
	Host    string `toml:"host" json:"host"`
	Port    int    `toml:"port" json:"port"`
	Listen  int    `toml:"listen" json:"listen"`
	Chanels int    `toml:"chanels" json:"chanels"`
	Diaps   int    `toml:"diaps" json:"diaps"`
	Diap    int    `toml:"diap" json:"diap"`
}
type Modbus struct {
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
	UId      int    `toml:"uid" json:"uid"`
	Debug    bool   `toml:"debug" json:"debug"`
	Log      bool   `toml:"log" json:"log"`
	Tmin     int    `toml:"tmin" json:"min"` //Минимальная длительность фазы
}

func (s *Setup) Update(es ExtSetup) {
	s.Utopia = es.Utopia
	s.Modbus = es.Modbus
	s.STCIP = es.STCIP
	s.ModbusRadar = es.ModbusRadar
	s.TrafficData = es.TrafficData
}
func (es *ExtSetup) Update(s Setup) {
	es.Utopia = s.Utopia
	es.Modbus = s.Modbus
	es.STCIP = s.STCIP
	es.ModbusRadar = s.ModbusRadar
	es.TrafficData = s.TrafficData
}
