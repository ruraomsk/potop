package traffic

import (
	"bufio"
	"encoding/xml"
	"net"

	"github.com/ruraomsk/ag-server/logger"
)

func serverWorker(socket net.Conn) {
	defer socket.Close()
	var sendto chan string
	var subscription chan bool

	sendto = make(chan string)
	subscription = make(chan bool)
	go serverSender(socket, sendto, subscription)
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
		// logger.Debug.Printf("Message %v\n", mess.Type)
		switch mess.Type {
		case "GetConnectionInformation":
			sendto <- GetConnectionInformation(buffer)
		case "ProductInformation":
			sendto <- GetProductInformation(buffer)
		case "GetTime":
			sendto <- GetTime(buffer)
		case "GetUpTime":
			sendto <- GetUpTime(buffer)
		case "GetImage":
			sendto <- GetImageMessage(buffer)
		case "Subscription":
			sendto <- GetSubcription(buffer, subscription)
		default:
			logger.Error.Println("Не распознано")
			sendto <- errorMessage
		}
	}
}
