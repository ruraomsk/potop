package web

import (
	"strconv"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/setup"
)

const setupText = `
		ListLayout {
			style = showPage,
			orientation = vertical,
			content = [
				TextView {
					style=header1,
					text = "<b>Изменение настроек связи с контроллером </b>",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						TextView {
							text = "Устройство",text-size="24px",
						},
						EditView{
							id=idDevice,type=text
						},
						TextView {
							text = "Uid",text-size="24px",
						},
						NumberPicker {
							id=idUid,type=editor,min=0,max=255,value=247
						},
					]
				},
				TextView {
					text = "<b>Изменение настроек связи со СПОТ Utopia </b>",text-align="center",text-size="24px",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						TextView {
							text = "Устройство",text-size="24px",
						},
						EditView{
							id=idUtopia,type=text
						},
						TextView {
							text = "Скорость ",text-size="24px",
						},
						NumberPicker {
							id=idBaudrate,type=editor,min=0,max=115200,value=38400
						},
					]
				},
				TextView {
					text = "<b>Изменение настроек связи с радарным комплексом</b>",text-align="center",text-size="24px",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						DropDownList { id =listRadar, current = 0, items = ["Включен", "Отключен"]},
						TextView {
							text = "IP Host",text-size="24px",
						},
						EditView{
							id=idIPRadar,type=text
						},
						TextView {
							text = "Номер порта",text-size="24px",
						},
						NumberPicker {
							id=idPortRadar,type=editor,min=0,max=32000,value=15001
						},
						TextView {
							text = "UID",text-size="24px",
						},
						NumberPicker {
							id=idUidRadar,type=editor,min=0,max=255,value=11
						},
						TextView {
							text = "Каналов",text-size="24px",
						},
						NumberPicker {
							id=idChanelsRadar,type=editor,min=0,max=16,value=16
						},
					]
				},
				TextView {
					text = "<b>Изменение настроек связи с TrafficData</b>",text-align="center",text-size="24px",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						DropDownList { id =listTrafficData, current = 0, items = ["Включен", "Отключен"]},
						TextView {
							text = "Server Host",text-size="24px",
						},
						EditView{
							id=idIPTrafficData,type=text
						},
						TextView {
							text = "Номер порта",text-size="24px",
						},
						NumberPicker {
							id=idPortTrafficData,type=editor,min=0,max=32000,value=0
						},
						TextView {
							text = "Listen port",text-size="24px",
						},
						NumberPicker {
							id=idListenTrafficData,type=editor,min=0,max=32000,value=0
						},
						TextView {
							text = "Каналов",text-size="24px",
						},
						NumberPicker {
							id=idChanelsTrafficData,type=editor,min=0,max=16,value=0
						},
					]
				},
				Button {
					id=setUpdate,content="Применить изменения"
				},
				TextView {
					id=idError,text = "",text-align="center",text-size="24px",text-color="red",
				},

			]
		}
`

func setupShow(session rui.Session) rui.View {

	view := rui.CreateViewFromText(session, setupText)
	if view == nil {
		return nil
	}

	rui.Set(view, "idIPRadar", "text", setup.ExtSet.ModbusRadar.Host)
	rui.Set(view, "idPortRadar", "value", setup.ExtSet.ModbusRadar.Port)
	rui.Set(view, "idUidRadar", "value", setup.ExtSet.ModbusRadar.ID)
	rui.Set(view, "idChanelsRadar", "value", setup.ExtSet.ModbusRadar.Chanels)
	c := 0
	if !setup.ExtSet.ModbusRadar.Radar {
		c = 1
	}
	rui.Set(view, "listRadar", "current", c)
	rui.Set(view, "listRadar", rui.DropDownEvent, func(_ rui.DropDownList, number int) {
		switch number {
		case 0:
			setup.ExtSet.ModbusRadar.Radar = true

		case 1:
			setup.ExtSet.ModbusRadar.Radar = false
		}
		if setup.ExtSet.TrafficData.Work && setup.ExtSet.ModbusRadar.Radar {
			rui.Set(view, "idError", "text", "<b>Нельзя одновременно указывать радар и TrafficData</b>")
		} else {
			rui.Set(view, "idError", "text", "")
		}
	})
	c = 0
	if !setup.ExtSet.TrafficData.Work {
		c = 1
	}
	rui.Set(view, "listTrafficData", "current", c)
	rui.Set(view, "listTrafficData", rui.DropDownEvent, func(_ rui.DropDownList, number int) {
		switch number {
		case 0:
			setup.ExtSet.TrafficData.Work = true

		case 1:
			setup.ExtSet.TrafficData.Work = false
		}
		if setup.ExtSet.TrafficData.Work && setup.ExtSet.ModbusRadar.Radar {
			rui.Set(view, "idError", "text", "<b>Нельзя одновременно указывать радар и TrafficData</b>")
		} else {
			rui.Set(view, "idError", "text", "")
		}
	})

	rui.Set(view, "idIPTrafficData", "text", setup.ExtSet.TrafficData.Host)
	rui.Set(view, "idPortTrafficData", "value", setup.ExtSet.TrafficData.Port)
	rui.Set(view, "idListenTrafficData", "value", setup.ExtSet.TrafficData.Listen)
	rui.Set(view, "idChanelsTrafficData", "value", setup.ExtSet.TrafficData.Chanels)

	rui.Set(view, "idDevice", "text", setup.ExtSet.Modbus.Device)
	rui.Set(view, "idUid", "value", setup.ExtSet.Modbus.UId)
	rui.Set(view, "idUtopia", "text", setup.ExtSet.Utopia.Device)
	rui.Set(view, "idBaudrate", "value", setup.ExtSet.Utopia.BaudRate)

	rui.Set(view, "setUpdate", rui.ClickEvent, func(rui.View) {

		setup.ExtSet.ModbusRadar.Host = rui.GetText(view, "idIPRadar")
		setup.ExtSet.ModbusRadar.Port = getInteger(rui.Get(view, "idPortRadar", "value"))
		setup.ExtSet.ModbusRadar.ID = getInteger(rui.Get(view, "idUidRadar", "value"))
		setup.ExtSet.ModbusRadar.Chanels = getInteger(rui.Get(view, "idChanelsRadar", "value"))
		logger.Info.Printf("Изменили радар на %s:%d uid %d каналов %d",
			setup.ExtSet.ModbusRadar.Host, setup.ExtSet.ModbusRadar.Port, setup.ExtSet.ModbusRadar.ID,
			setup.ExtSet.ModbusRadar.Chanels)

		setup.ExtSet.TrafficData.Host = rui.GetText(view, "idIPTrafficData")
		setup.ExtSet.TrafficData.Port = getInteger(rui.Get(view, "idPortTrafficData", "value"))
		setup.ExtSet.TrafficData.Listen = getInteger(rui.Get(view, "idListenTrafficData", "value"))
		setup.ExtSet.TrafficData.Chanels = getInteger(rui.Get(view, "idChanelsTrafficData", "value"))
		logger.Info.Printf("Изменили TrafficData на %s:%d  %d каналов %d",
			setup.ExtSet.TrafficData.Host, setup.ExtSet.TrafficData.Port, setup.ExtSet.TrafficData.Listen,
			setup.ExtSet.TrafficData.Chanels)

		setup.ExtSet.Modbus.Device = rui.GetText(view, "idDevice")
		setup.ExtSet.Modbus.UId = getInteger(rui.Get(view, "idUid", "value"))
		logger.Info.Printf("Изменили КДМ на %s uid %d", setup.ExtSet.Modbus.Device, setup.ExtSet.Modbus.UId)
		setup.ExtSet.Utopia.Device = rui.GetText(view, "idUtopia")
		setup.ExtSet.Utopia.BaudRate = getInteger(rui.Get(view, "idBaudrate", "value"))
		logger.Info.Printf("Изменили Utopia на %s baudrate %d", setup.ExtSet.Utopia.Device, setup.ExtSet.Utopia.BaudRate)
		saveAndExit <- 1
	})

	return view
}
func getInteger(a any) (result int) {
	s, ok := a.(string)
	if ok {
		result, _ = strconv.Atoi(s)
	}
	f, ok := a.(float64)
	if ok {
		result = int(f)
	}
	return
}
