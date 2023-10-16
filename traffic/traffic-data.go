package traffic

import (
	"encoding/xml"
	"time"
)

type Message struct {
	Type string `xml:",attr"`
}
type ComnectionInformation struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Ip      string `xml:",attr"`
		Subnet  string `xml:",attr"`
		Gateway string `xml:",attr"`
	}
}
type Time struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Utc          string `xml:",attr"`
		Milliseconds string `xml:",attr"`
	}
}
type UpTime struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Value string `xml:",attr"`
	}
}
type GetImage struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		CameraId string `xml:",attr"`
	}
}
type Image struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		CameraId string `xml:",attr"`
		Width    string `xml:",attr"`
		Height   string `xml:",attr"`
		Format   string `xml:",attr"`
		Data     string `xml:",attr"`
	}
}
type Camera struct {
	Id           string `xml:",attr"`
	Name         string `xml:",attr"`
	Url          string `xml:",attr"`
	StreamingUrl string `xml:",attr"`
}
type ProductInformation struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Cameras string `xml:",attr"`
		Company string `xml:",attr"`
		Product string `xml:",attr"`
		Version string `xml:",attr"`
		Camera  []Camera
	}
}
type Subs struct {
	Type        string `xml:",attr"`
	Action      string `xml:",attr"`
	ReturnValue string `xml:",attr"`
}
type Subscription struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Subscription []Subs `xml:"Subscription"`
	}
}

type CarMoved struct {
	XMLName xml.Name `xml:"Message"`
	Type    string   `xml:",attr"`
	Body    struct {
		Type       string `xml:",attr"`
		ZoneId     string `xml:",attr"`
		CarId      string `xml:",attr"`
		CarType    string `xml:",attr"`
		Utc        string `xml:",attr"`
		TimeInZone string `xml:",attr"`
	}
}

var startTime time.Time

const countCameras = 16
const endline byte = '\n'
const errorMessage = `
<Message Type="ErrorMessage" />
`
