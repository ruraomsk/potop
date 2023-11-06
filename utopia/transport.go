package utopia

import (
	"errors"
	"sync"
	"time"

	"github.com/goburrow/serial"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
)

type StatusUtopiaTransport struct {
	mutex         sync.Mutex
	Connect       bool
	LastOperation time.Time
	FromServer    []byte
	ToServer      []byte
}

func (s *StatusUtopiaTransport) getConnect() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.Connect
}
func (s *StatusUtopiaTransport) setConnect(set bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Connect = set
}
func (s *StatusUtopiaTransport) setFromServer(set []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.FromServer = set
	s.LastOperation = time.Now()
}
func (s *StatusUtopiaTransport) setToServer(set []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.ToServer = set
	s.LastOperation = time.Now()
}

var fromServer chan []byte
var toServer chan []byte
var fromController chan []byte
var toController chan []byte
var port serial.Port
var err error
var statusTransport = StatusUtopiaTransport{Connect: false, LastOperation: time.Unix(0, 0), FromServer: make([]byte, 0), ToServer: make([]byte, 0)}

func Transport() {

	fromController = make(chan []byte)
	fromServer = make(chan []byte)
	toController = make(chan []byte)
	toServer = make(chan []byte)
	if setup.Set.Utopia.Debug {
		statusTransport.setConnect(true)
		for {
			u := <-toController
			statusTransport.setFromServer(u)
			fromServer <- u
			u = <-toServer
			statusTransport.setToServer(u)
			fromController <- u
		}
	}
	count := 0
	config := serial.Config{Address: setup.Set.Utopia.Device, BaudRate: setup.Set.Utopia.BaudRate, StopBits: 0, Parity: "N", Timeout: 5 * time.Second}
	for {
		if !statusTransport.getConnect() {
			time.Sleep(5 * time.Second)
			port, err = serial.Open(&config)
			if err != nil {
				if (count % 10) == 0 {
					logger.Error.Printf("spot open port %s %d %v", setup.Set.Utopia.Device, setup.Set.Utopia.BaudRate, err.Error())
				}
				count++
				continue
			}
			statusTransport.setConnect(true)
			count = 0
		}
		logger.Info.Printf("spot open port %s %d ", setup.Set.Utopia.Device, setup.Set.Utopia.BaudRate)
		for statusTransport.getConnect() {
			buffer, err := getFromServer()
			if err != nil {
				logger.Error.Printf("recieve from spot %s", err.Error())
				port.Close()
				statusTransport.setConnect(false)
				continue
			}
			fromServer <- buffer
			statusTransport.setFromServer(buffer)
			buffer = <-toServer
			err = sendToServer(buffer)
			if err != nil {
				logger.Error.Printf("send to spot %s", err.Error())
				port.Close()
				statusTransport.setConnect(false)
				continue
			}
			statusTransport.setToServer(buffer)
		}
	}
}
func getFromServer() ([]byte, error) {
	body := make([]byte, 256)
	n, err := port.Read(body)
	if err != nil {
		return body, err
	}
	// logger.Debug.Printf("from Utopia % 02X", body[:n])
	return body[:n], nil
}
func sendToServer(buffer []byte) error {
	n, err := port.Write(buffer)
	if err != nil {
		return err
	}
	if n != len(buffer) {
		return errors.New("отправлен не весь буфер")
	}
	return nil
}
