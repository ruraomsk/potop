package traffic

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetConnectionInformation(message string) string {
	ci := ComnectionInformation{Type: "ConnectionInformation"}
	ci.Body.Ip = "192.168.100.2"
	ci.Body.Subnet = "255.255.255.0"
	ci.Body.Gateway = "127.0.0.1"
	s, err := xml.Marshal(ci)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
func GetProductInformation(message string) string {
	pi := ProductInformation{Type: "ProductInformation"}
	pi.Body.Cameras = fmt.Sprintf("%d", countCameras)
	pi.Body.Company = "TrassicData"
	pi.Body.Product = "TD TraffiCam x-stream"
	pi.Body.Version = "V1.0"
	pi.Body.Camera = make([]Camera, 0)
	for i := 1; i <= countCameras; i++ {
		c := Camera{Id: fmt.Sprintf("%d", i),
			Name:         fmt.Sprintf("Camera %d", i),
			Url:          fmt.Sprintf("rtsp://192.168.100.%d:%d/%d", rand.Intn(100), rand.Intn(9999), rand.Intn(200)),
			StreamingUrl: fmt.Sprintf("rtsp://192.168.100.%d:%d/edge%d.stream", rand.Intn(100), rand.Intn(9999), rand.Intn(200))}
		pi.Body.Camera = append(pi.Body.Camera, c)
	}
	s, err := xml.Marshal(pi)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
func GetTime(message string) string {
	ti := Time{Type: "Time"}
	ti.Body.Utc = fmt.Sprintf("%d", time.Now().Unix())
	ti.Body.Milliseconds = "0"
	s, err := xml.Marshal(ti)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
func GetUpTime(message string) string {
	uti := UpTime{Type: "UpTime"}
	uti.Body.Value = fmt.Sprintf("%d", int64(time.Now().Sub(startTime).Seconds()))
	s, err := xml.Marshal(uti)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
func GetImageMessage(message string) string {
	var gi GetImage
	err := xml.Unmarshal([]byte(message), &gi)
	if err != nil {
		return errorMessage
	}
	camera, _ := strconv.Atoi(gi.Body.CameraId)
	if camera < 1 || camera > countCameras {
		return errorMessage
	}
	image := Image{Type: "Image"}
	image.Body.CameraId = gi.Body.CameraId
	image.Body.Width = "640"
	image.Body.Height = "420"
	image.Body.Format = "JPEG"
	image.Body.Data = "[encoded image]"
	s, err := xml.Marshal(image)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
func GetSubcription(message string, subs chan bool) string {
	var sb Subscription
	err := xml.Unmarshal([]byte(message), &sb)
	if err != nil {
		return errorMessage
	}
	// logger.Debug.Printf("in subscription %v", sb)
	ok := "BAD"
	for _, v := range sb.Body.Subscription {
		if strings.Compare(v.Type, "CarExit") == 0 {
			ok = "OK"
			if strings.Compare(v.Action, "Subscribe") == 0 {
				subs <- true
			} else {
				subs <- false
			}
		}
	}
	rep := Subscription{Type: "Subscription"}
	rep.Body.Subscription = make([]Subs, 0)
	rep.Body.Subscription = append(rep.Body.Subscription, Subs{Type: "CarExit", ReturnValue: ok})

	s, err := xml.Marshal(rep)
	if err != nil {
		return errorMessage
	}
	return string(s)
}
