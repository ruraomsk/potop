package traffic

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

type move struct {
	zoneid     int
	carid      int
	cartype    int
	timeinzone float32
}

func sendString(socket net.Conn, message string) error {
	buffer := make([]byte, 0)
	buffer = append(buffer, []byte(message)...)
	buffer = append(buffer, endline)
	_, err := socket.Write(buffer)
	return err
}
func serverSender(socket net.Conn, sendto chan string, subscription chan bool) {
	connected := false
	var Chanels bool
	var soc net.Conn
	var err error
	host := strings.Split(socket.RemoteAddr().String(), ":")
	for {
		for !connected {
			soc, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host[0], 1667))
			if err != nil {
				logger.Error.Println(err.Error())
				slp := time.NewTimer(10 * time.Second)
				to := true
				for to {
					select {
					case s := <-sendto:
						logger.Error.Printf("Не отправлено %s", s)
					case Chanels = <-subscription:

					case <-slp.C:
						to = false
					}
				}
			} else {
				connected = true
			}
		}
		logger.Info.Printf("Начинаем работу с %s", socket.RemoteAddr())
		tick := time.NewTicker(250 * time.Millisecond)
		for {
			select {
			case s := <-sendto:

				err = sendString(soc, s)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

			case Chanels = <-subscription:
				// logger.Debug.Printf("Подписка %v", Chanels)
			case <-tick.C:
				if !Chanels {
					continue
				}
				// logger.Debug.Println("Строим ...")
				moved := make([]move, 0)
				for i := 1; i <= countCameras; i++ {
					for j := 0; j < rand.Intn(3); j++ {
						if rand.Intn(2) == 1 {
							moved = append(moved, move{zoneid: i, carid: rand.Intn(10000) + 1, cartype: rand.Intn(15) + 1,
								timeinzone: rand.Float32() * 20})
						}
					}
				}
				// logger.Debug.Printf("Набрали %d", len(moved))

				for _, v := range moved {
					mm := CarMoved{Type: "Event"}
					mm.Body.Type = "CarExit"
					mm.Body.CarId = fmt.Sprintf("%d", v.carid)
					mm.Body.CarType = fmt.Sprintf("%d", v.cartype)
					mm.Body.ZoneId = fmt.Sprintf("%d", v.zoneid)
					mm.Body.TimeInZone = fmt.Sprintf("%f", v.timeinzone)
					mm.Body.Utc = fmt.Sprintf("%d", time.Now().Unix())
					ss, err := xml.Marshal(mm)
					if err != nil {
						logger.Error.Println(err.Error())
						continue
					}
					err = sendString(soc, string(ss))
					if err != nil {
						logger.Error.Println(err.Error())
						break
					}

				}
			}

		}
	}
}
