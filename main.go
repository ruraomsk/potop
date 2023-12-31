package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/anoshenko/rui"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/diagramm"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/radar"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/stat"
	"github.com/ruraomsk/potop/stcip"
	"github.com/ruraomsk/potop/traffic"
	"github.com/ruraomsk/potop/utopia"
	"github.com/ruraomsk/potop/web"
)

func init() {
	setup.Set = new(setup.Setup)
	setup.ExtSet = new(setup.ExtSetup)
	if _, err := toml.DecodeFS(resources, "config/base.toml", &setup.Set); err != nil {
		fmt.Println("Dismiss base.toml")
		os.Exit(-1)
		return
	}
	if _, err := os.Stat("config.json"); err == nil {
		file, err := os.ReadFile("config.json")
		if err == nil {
			_ = json.Unmarshal(file, &setup.ExtSet)
			setup.Set.Update(*setup.ExtSet)
		}
	}
	setup.ExtSet.Update(*setup.Set)
	_ = os.MkdirAll(setup.Set.LogPath, 0777)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	fmt.Println("Potop start")
	logger.Info.Println("Potop start")
	go journal.LoggerStart()
	go hardware.Start()
	go diagramm.DiagrammStart()

	time.Sleep(time.Second)
	if setup.Set.Utopia.Run {
		go utopia.Transport()
		if setup.Set.Utopia.Debug {
			go utopia.Server()
		}
		go utopia.Controller()
	}
	if setup.Set.STCIP.Run {
		go stcip.Start()
	}
	isStat := false
	if setup.Set.ModbusRadar.Work {
		go stat.Start(setup.Set.ModbusRadar.Chanels, setup.Set.ModbusRadar.Diaps)
		go radar.Radar(setup.Set.ModbusRadar.Diap)
		isStat = true
	}
	if setup.Set.TrafficData.Work {
		go stat.Start(setup.Set.TrafficData.Chanels, setup.Set.TrafficData.Diaps)
		go traffic.Start(setup.Set.TrafficData.Diap)
		isStat = true
	}
	if !isStat {
		go stat.NoStatistics()
	}

	go rui.AddEmbedResources(&resources)
	go web.Web()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nwait ...")
	time.Sleep(1 * time.Second)
	fmt.Println("Potop stop")
	logger.Info.Println("Potop stop")
}
