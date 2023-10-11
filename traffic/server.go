package traffic

import (
	"fmt"
	"net"
	"time"

	"github.com/ruraomsk/ag-server/logger"
)

func Server(port int) {
	startTime = time.Now()
	logger.Info.Printf("Слушаем %d", port)
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
		go serverWorker(socket)
	}
}
