package traffic

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
)

type StatusTrafficData struct {
	StatusConnectionInformation ComnectionInformation
	StatusProductInformation    ProductInformation
	StatusGetTime               Time
	StatusUpTime                UpTime
	StatusNowSubscription       Subscription
	StatusWorkClient            bool
	StatusWorkReciever          bool
	Datas                       []int
	LastOperation               time.Time
}

var socketClient net.Conn
var connectionInformation ComnectionInformation
var productInformation ProductInformation
var getTime Time
var upTime UpTime
var image Image
var nowSubscription Subscription

var mutex sync.Mutex
var datas []int
var lastDatas []int
var workClient = false
var workReciever = false
var lastOperation time.Time = time.Unix(0, 0)

func GetStatusTrafficData() StatusTrafficData {
	mutex.Lock()
	defer mutex.Unlock()
	return StatusTrafficData{
		StatusConnectionInformation: connectionInformation,
		StatusProductInformation:    productInformation,
		StatusGetTime:               getTime,
		StatusUpTime:                upTime,
		StatusNowSubscription:       nowSubscription,
		StatusWorkClient:            workClient,
		StatusWorkReciever:          workReciever,
		Datas:                       lastDatas,
		LastOperation:               lastOperation,
	}
}
func GetDatas() ([]int, time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	return lastDatas, lastOperation
}
func GetValues() string {
	if !workClient || !workReciever {
		return "Нет связи с сервером TrafficData"
	}
	mutex.Lock()
	defer mutex.Unlock()
	return fmt.Sprintf("%s %v", lastOperation.Format("15:04:05"), lastDatas)
}
func GetStatus() string {
	return fmt.Sprintf("Server %s:%d", setup.Set.TrafficData.Host, setup.Set.TrafficData.Port)
}

func Start() {
	if !setup.Set.TrafficData.Work {
		return
	}
	if setup.Set.TrafficData.Debug {
		go Server(setup.Set.TrafficData.Port)
	}
	go Client(setup.Set.TrafficData.Host, setup.Set.TrafficData.Port, setup.Set.TrafficData.Listen)
}
func Client(host string, port int, listen int) {
	go clientReciever(listen)
	var err error
	datas = make([]int, setup.Set.TrafficData.Chanels)
	lastDatas = make([]int, len(datas))
	for i := 0; i < len(datas); i++ {
		datas[i] = 0
	}
	for {
		workClient = false
		for !workClient {
			socketClient, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				logger.Error.Println(err.Error())
				time.Sleep(10 * time.Second)
			} else {
				workClient = true
			}
		}

		logger.Info.Printf("TrafficData начинаем работу с %s", host)
		for i := 0; i < 1; i++ {
			err = sendConnectionInformation(socketClient)
			if err != nil {
				break
			}
			err = sendProductInformation(socketClient)
			if err != nil {
				break
			}
			err = sendSubscription(socketClient)
			if err != nil {
				break
			}
			err = sendGetTime(socketClient)
			if err != nil {
				break
			}
			err = sendUpTime(socketClient)
			if err != nil {
				break
			}
		}
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		ticker := time.NewTicker(1 * time.Second)
		oneSecond := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				err = sendGetTime(socketClient)
				if err != nil {
					break
				}
				err = sendUpTime(socketClient)
				if err != nil {
					break
				}
			case <-oneSecond.C:
				saveDatas()
			}
		}
	}
}
func saveDatas() {
	mutex.Lock()
	for i := 0; i < len(datas); i++ {
		lastDatas[i] = datas[i]
		datas[i] = 0
	}
	var send []OneTick
	if time.Now().Sub(lastOperation).Seconds() > 2 {
		send = badStatistics()
	} else {
		send = goodStatistics()
	}
	mutex.Unlock()
	for _, v := range send {
		logger.Info.Printf("%v", v)
	}

}
func badStatistics() []OneTick {
	r := make([]OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.TrafficData.Chanels; i++ {
		r = append(r, OneTick{Nomber: i + 3, Value: Value{Status: 1, Time: t, Value: 0}})
	}
	return r
}
func goodStatistics() []OneTick {
	r := make([]OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.TrafficData.Chanels; i++ {
		r = append(r, OneTick{Nomber: i + 3, Value: Value{Status: 0, Time: t, Value: int(lastDatas[i])}})
	}
	return r
}
