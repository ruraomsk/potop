package transport

import (
	"fmt"
	"time"
)

var CommandOids = []Define{
	{Code: CodeCallPhase, OID: "1.3.6.1.4.1.1618.3.7.2.11.1"},
	{Code: CodeCallPlan, OID: "1.3.6.1.4.1.1618.3.7.2.2.1"},
	{Code: CodeCallFlash, OID: "1.3.6.1.4.1.1618.3.2.2.1.1"},
	{Code: CodeCallDark, OID: "1.3.6.1.4.1.1618.3.2.2.2.1"},
	{Code: CodeCallAllRed, OID: "1.3.6.1.4.1.1618.3.2.2.2.1"},
}

func receiverCommands() {

	for {
		code := <-CommandFromWeb
		for _, v := range CommandOids {
			if code.Code == v.Code {
				toHistory <- History{Time: time.Now(), Type: 1, Message: fmt.Sprintf("%s %d", v.OID, code.Value)}
				Commander <- Command{OID: v.OID, Code: v.Code, Value: code.Value}
				break
			}
		}
	}
}
