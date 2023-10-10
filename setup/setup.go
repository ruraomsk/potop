package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath string `toml:"logpath"`
	Utopia  Utopia `toml:"utopia"`
}
type Utopia struct {
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
	UId      int    `toml:"uid" json:"uid"`
	Work     bool   `toml:"work" json:"work"`
	Debug    bool   `toml:"debug" json:"debug"`
}
