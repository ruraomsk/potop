package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath     string      `toml:"logpath"`
	Utopia      Utopia      `toml:"utopia"`
	TrafficData TrafficData `toml:"trafficdata"`
}
type ExtSetup struct {
	Utopia      Utopia      `toml:"utopia" json:"utopia"`
	TrafficData TrafficData `toml:"trafficdata" json:"trafficdata"`
}

type Utopia struct {
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
	UId      int    `toml:"uid" json:"uid"`
	Work     bool   `toml:"work" json:"work"`
	Debug    bool   `toml:"debug" json:"debug"`
}
type TrafficData struct {
	Work    bool   `toml:"work" json:"work"`
	Debug   bool   `toml:"debug" json:"debug"`
	Host    string `toml:"host" json:"host"`
	Port    int    `toml:"port" json:"port"`
	Listen  int    `toml:"listen" json:"listen"`
	Chanels int    `toml:"chanels" json:"chanels"`
}

func (s *Setup) Update(es ExtSetup) {
	s.Utopia = es.Utopia
	// s.Modbus = es.Modbus
	// s.VisioDevice = es.VisioDevice
	// s.ModbusRadar = es.ModbusRadar
	s.TrafficData = es.TrafficData
}
func (es *ExtSetup) Update(s Setup) {
	es.Utopia = s.Utopia
	// es.Modbus = s.Modbus
	// es.VisioDevice = s.VisioDevice
	// es.ModbusRadar = s.ModbusRadar
	es.TrafficData = s.TrafficData
}
