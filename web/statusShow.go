package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"

	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/radar"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/stcip"
	"github.com/ruraomsk/potop/traffic"
	"github.com/ruraomsk/potop/utopia"
)

// border = _{ style = solid, width = 1px, color = darkgray },
const statusText = `
		ListLayout {
			style = showPage,
			orientation = vertical,
			padding="16px",
			content = [
				TextView {
					style=header1,
					id=titleStatus,text = "",
				},
				TextView {
					id=idUtopia,
					text = "",
					text-size="24px",
				},
				TextView {
					id=idSTCIP,
					text = "",
					text-size="24px",
				},
				TextView {
					id=idModbus,
					text = "",
					text-size="24px",
				},
				TextView {
					id=setModbusRadar,
					text = "",
					text-size="24px",
				},
				TextView {
					id=setTrafficData,
					text = "",
					text-size="24px",
				},
			]
		}
`

func toString(t time.Time) string {
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func makeViewStatus(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	t := time.Now()
	rui.Set(view, "titleStatus", "text", fmt.Sprintf("<b>Текущее состояние контроллера %d </b>%02d:%02d:%02d", setup.Set.Id,
		t.Hour(), t.Minute(), t.Second()))
	c := ""
	if setup.Set.Utopia.Run {
		c = fmt.Sprintf("Соединение со СПОТ Utopia device %s baud %d \t",
			setup.Set.Utopia.Device, setup.Set.Utopia.BaudRate)
		if utopia.GetConnect() {
			c += "установлено"
		} else {
			c += "отсутствует"
		}
	} else {
		c = "Utopia интерфейс отключен"
	}
	rui.Set(view, "idUtopia", "text", c)
	if setup.Set.STCIP.Run {
		c = fmt.Sprintf("Соединение со STCIP device %s port %d listen %d \t",
			setup.Set.STCIP.Host, setup.Set.STCIP.Port, setup.Set.STCIP.Listen)
		if stcip.GetConnect() {
			c += "установлено"
		} else {
			c += "отсутствует"
		}
	} else {
		c = "STCIP интерфейс отключен"
	}

	rui.Set(view, "idSTCIP", "text", c)

	c = fmt.Sprintf("Соединение с КДМ device %s baud %d parity %s uid %d \t",
		setup.Set.Modbus.Device, setup.Set.Modbus.BaudRate, setup.Set.Modbus.Parity, setup.Set.Modbus.UId)
	if hardware.StateHardware.GetConnect() {
		c += "установлено"
	} else {
		c += "отсутствует"
	}
	rui.Set(view, "idModbus", "text", c)
	if !setup.Set.ModbusRadar.Work {
		rui.Set(view, "setModbusRadar", "text", "Оключен прием данных от радаров")
	} else {
		rui.Set(view, "setModbusRadar", "text", fmt.Sprintf("От радаров (%s): %s ", radar.GetStatus(), radar.GetValues()))
	}
	if !setup.Set.TrafficData.Work {
		rui.Set(view, "setTrafficData", "text", "Оключен прием данных от TrafficData")
	} else {
		rui.Set(view, "setTrafficData", "text", fmt.Sprintf("От TrafficData (%s): %s ", traffic.GetStatus(), traffic.GetValues()))
	}
}
func updaterStatus(view rui.View, session rui.Session) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if view == nil {
			return
		}
		w, ok := SessionStatus[session.ID()]
		if !ok {
			continue
		}

		if !w {
			continue
		}
		makeViewStatus(view)
	}
}

func statusShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, statusText)
	if view == nil {
		return nil
	}
	makeViewStatus(view)
	go updaterStatus(view, session)

	return view
}
