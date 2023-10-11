package traffic

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

func clientReciever(port int) {
	logger.Info.Printf("Ждем сервер на %d", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		go clientWorker(socket)
	}
}
func clientWorker(socket net.Conn) {
	defer func() {
		socket.Close()
		workReciever = false
	}()
	workReciever = true
	reader := bufio.NewReader(socket)
	for {
		buffer, err := reader.ReadString(endline)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		var mess Message
		err = xml.Unmarshal([]byte(buffer), &mess)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		lastOperation = time.Now()
		// logger.Debug.Printf("Message from server %v\n", mess.Type)
		switch mess.Type {
		case "ConnectionInformation":
			replayConnectionInformation(buffer)
		case "ProductInformation":
			replayProductInformation(buffer)
		case "Time":
			replayTime(buffer)
		case "UpTime":
			replayUpTime(buffer)
		case "Image":
			replayImage(buffer)
		case "Subscription":
			replaySubcription(buffer)
		case "Event":
			replayEvent(buffer)
		default:
			logger.Error.Println("Не распознано")
		}
	}

}
