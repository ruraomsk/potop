package web

import (
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/hardware"
)

const controlText = `
		ListLayout {
			style = showPage,
			orientation = vertical,
			content = [
				TextView {
					style=header1,
					text = "<b>Текущее состояние контроллера</b>",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						TextView {
							style=header1,
							id=idAllRed,text="Кругом Красный"
						},
						TextView {
							style=header1,
							id=idFlashing,text="Желтое Мигание"
						},
						TextView {
							style=header1,
							id=idDark,text="Выключено"
						},
					]
				},
				TextView {
					style=header1,
					text = "<b>Изменение режима работы контроллера </b>",
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						Button {
							id=setAllRed,content="Установить Кругом Красный"
						},
						Button {
							id=setFlashing,content="Установить Желтое Мигание"
						},
						Button {
							id=setDark,content="Выключить"
						},
						Button {
							id=setLocal,content="Вернуть в ЛР"
						},
					]
				},
			]
		}
`

func makeViewControl(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	hs := hardware.GetStateHard()
	if hs.Dark {
		rui.Set(view, "idDark", "visibility", "visible")
	} else {
		rui.Set(view, "idDark", "visibility", "invisible")
	}
	if hs.AllRed {
		rui.Set(view, "idAllRed", "visibility", "visible")
	} else {
		rui.Set(view, "idAllRed", "visibility", "invisible")
	}
	if hs.Flashing {
		rui.Set(view, "idFlashing", "visibility", "visible")
	} else {
		rui.Set(view, "idFlashing", "visibility", "invisible")
	}
}
func updaterControl(view rui.View, session rui.Session) {
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
		makeViewControl(view)
	}

}
func controlKDM(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, controlText)
	if view == nil {
		return nil
	}
	rui.Set(view, "setAllRed", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(4, 0)
		logger.Info.Printf("Оператор установил КК")
	})
	rui.Set(view, "setFlashing", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(3, 0)
		logger.Info.Printf("Оператор установил ЖМ")
	})
	rui.Set(view, "setDark", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(6, 0)
		logger.Info.Printf("Оператор установил ОС")
	})
	rui.Set(view, "setLocal", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(1, 0)
		logger.Info.Printf("Оператор установил ЛР")
	})
	makeViewControl(view)
	go updaterControl(view, session)
	return view
}
