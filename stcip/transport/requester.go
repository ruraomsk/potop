package transport

import "time"

var ReqOids = []Define{
	{Code: CodeReqStatus, OID: "1.3.6.1.4.1.1618.3.6.2.1.2"},
	{Code: CodeReqMajor, OID: "1.3.6.1.4.1.1618.3.6.2.2.2"},
	{Code: CodeReqPlan, OID: "1.3.6.1.4.1.1618.3.5.2.1.7"},
	{Code: CodeReqSource, OID: "1.3.6.1.4.1.1618.3.7.2.1.3"},
	{Code: CodeReqPhase, OID: "1.3.6.1.4.1.1618.3.7.2.11.2"},
	{Code: CodeReqSignalGroups, OID: "1.3.6.1.4.1.1618.3.5.2.1.6"},
	{Code: CodeReqAlarm, OID: "1.3.6.1.4.1.1618.3.1.2.2.2"},
}

func receiveString(code int) string {
	for _, v := range ReqOids {
		if code == v.Code {
			return v.OID
		}
	}
	return "undef"
}
func receiverRequests() {
	for {
		for _, v := range ReqOids {
			time.Sleep(5 * time.Second)
			Requester <- Request{Code: v.Code, OID: v.OID}
		}
	}
}
