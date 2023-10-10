package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/utopia"
)

//go:embed config
var resources embed.FS

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(resources, "config/base.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis base.toml")
		os.Exit(-1)
		return
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	fmt.Println("Potop start")
	logger.Info.Println("Potop start")
	go utopia.Transport()
	time.Sleep(time.Second)
	go utopia.Server()
	go utopia.Controller()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nwait ...")
	time.Sleep(1 * time.Second)
	fmt.Println("Potop stop")
	logger.Info.Println("Potop stop")
}
