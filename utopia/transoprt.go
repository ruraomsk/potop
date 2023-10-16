package utopia

import (
	"fmt"

	"github.com/ruraomsk/potop/setup"
)

var fromServer chan []byte
var toServer chan []byte
var fromController chan []byte
var toController chan []byte

func Transport() {

	fromController = make(chan []byte, 100)
	fromServer = make(chan []byte, 100)
	toController = make(chan []byte, 100)
	toServer = make(chan []byte, 100)
	if setup.Set.Debug {
		for {
			select {
			case u := <-toController:
				fromServer <- u
			case u := <-toServer:
				fromController <- u
			}
		}
	} else {
		fmt.Println("Нет еще обмена с Utopia сервером!")
		panic(-1)
	}
}
