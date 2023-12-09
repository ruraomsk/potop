package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/radar"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/traffic"
)

const mainText = `
ListLayout {
	style = showPage,
	orientation = vertical,
	content = [
		GridLayout {
			style = optionsTable,
			border = _{ style = solid, width = 1px, color = darkgray },
			content = [
				TextView {
					row = 0, column = 0,
					id="isKDM",text-align = center,
					text = "КДМ",background-color=gray,
				},
				TextView {
					row = 0, column = 1,
					id="isCentral",text-align = center,
					text = "ЦЕНТР",background-color=gray,
				},
				TextView {
					row = 0, column = 2,
					id="isRadar",text-align = center,
					text = "РАДАРЫ",background-color=gray,
				},
				TextView {
					row = 0, column = 3,
					id="isTrafficData",text-align = center,
					text = "ТРАФИК ДАТА",background-color=gray,
				},
				TextView {
					row = 0, column = 4,
					id="time",
					text = "Текушее время",color=red,
				},

			]
		},
		GridLayout {
			content = [
				ListLayout {
					orientation = vertical,
					row = 0, column = 0,
					content = [
						TextView {
							id=idLine1,text = "Line1",text-size="18px",
						},
						TextView {
							id=idLine2,text = "Line2",text-size="18px",
						},
						TextView {
							id=idLine3,text = "Line3",text-size="18px",
						},
						TableView {
							style=table,
							cell-horizontal-align = left,  title = "Направления", id=idNaps,
						},
					]
				},
				GridLayout {
					style = optionsTable,
					border = _{ style = solid, width = 1px, color = darkgray },
					row = 0, column = 1,vertical-align = center,cell-vertical-align=center,
					content = [
						TextView {
							row = 0, column = 0,
							id="Flashing",text-align = center,
							text = "ЖМ",background-color=gray,
						},
						TextView {
							row = 0, column = 1,
							id="AllRed",text-align = center,
							text = "КК",background-color=gray,
						},
						TextView {
							row = 0, column = 2,
							id="Dark",text-align = center,
							text = "ОС",background-color=gray,
						},
						Button {
							row = 2, column = 1,
							id=setAllRedOn,content="КК вкл"
						},
						Button {
							row = 2, column = 0,
							id=setFlashingOn,content="ЖМ вкл"
						},
						Button {
							row = 2, column = 2,
							id=setDarkOn,content="ОС вкл"
						},
						Button {
							row = 5, column = 0:2,
							id=setLocal,content="Вернуть в ЛР"
						},
						Button {
							row = 6, column = 0:1,
							id=setPlan,content="Установить План"
						},
						NumberPicker {
							row = 6, column = 2,
							id=idPlan,type=editor,min=0,max=32,value=0
						},

					]
				},
			]
		},
		TableView {
			style=table,
			cell-horizontal-align = left,  title = "Диагностика", id=idDiags,
		},
	]
}
`

func updateMessages(view rui.View) {
	var content [][]any
	content = append(content, []any{"Время", "Сообщение"})
	count := 1
	ms := journal.GetMessages(1)
	for _, v := range ms {
		content = append(content, []any{toString(v.Time), v.Value})
		count++

	}
	rui.SetParams(view, "idDiags", rui.Params{
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}
func makeButtonOnScreen(view rui.View) {
	rui.Set(view, "setAllRedOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(4, 0)
		logger.Info.Printf("Оператор установил КК")
	})
	rui.Set(view, "setFlashingOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(3, 0)
		logger.Info.Printf("Оператор установил ЖМ")
	})
	rui.Set(view, "setDarkOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(6, 0)
		logger.Info.Printf("Оператор установил ОС")
	})
	rui.Set(view, "setLocal", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(1, 0)
		logger.Info.Printf("Оператор перевел в ЛР")
	})
	rui.Set(view, "setPlan", rui.ClickEvent, func(rui.View) {
		hardware.CommandUtopia(7, getInteger(rui.Get(view, "idPlan", "value")))
		logger.Info.Printf("Оператор вызвал план %d", getInteger(rui.Get(view, "idPlan", "value")))
	})

}
func updatePartKDM(view rui.View) {
	hs := hardware.GetStateHard()
	if hs.AllRed {
		rui.Set(view, "AllRed", "background-color", "green")
	} else {
		rui.Set(view, "AllRed", "background-color", "gray")
	}
	if hs.Flashing {
		rui.Set(view, "Flashing", "background-color", "green")
	} else {
		rui.Set(view, "Flashing", "background-color", "gray")
	}
	if hs.Dark {
		rui.Set(view, "Dark", "background-color", "green")
	} else {
		rui.Set(view, "Dark", "background-color", "gray")
	}

	source := "внутренний"
	if hs.SourceTOOB {
		source = "внешний"
	}
	rui.Set(view, "idLine1", "text", fmt.Sprintf("<b>Последняя команда в</b> %s <b>Тмин=%d Маска=%x остаток watchdog=%d</b>", toString(hs.LastOperation), hs.Tmin, hs.MaskCommand, hs.RealWatchDog))
	rui.Set(view, "idLine2", "text", fmt.Sprintf("<b>WatchDog</b> %d <b>План</b> %d <b>Статус КДМ</b> % 02X <b>Источник ТООВ</b> %s",
		hs.WatchDog, hs.Plan, hs.Status, source))
	rui.Set(view, "idLine3", "text", fmt.Sprintf("<b>Расшифровка статуса : %s </b>", hardware.GetError()))
	var content [][]any
	content = append(content, []any{"Нап", "Задание", "Состояние", "Счетчик ТООВ"})
	count := 1
	s := hs.MaskCommand
	for i := 0; i < 32; i++ {
		st := "Закрыто"
		if s&0x1 != 0 {
			st = "Открыто"
		}
		s = s >> 1
		ds := "undef"
		switch hs.StatusDirs[i] {
		case 0:
			ds = "все сигналы выключены"
		case 1:
			ds = "направление перешло в неактивное состояние, желтый после зеленого"
		case 2:
			ds = "направление перешло в неактивное состояние, красный"
		case 3:
			ds = "направление перешло в активное состояние, красный"
		case 4:
			ds = "направление перешло в активное состояние, красный c желтым"
		case 5:
			ds = "направление перешло в активное состояние, зеленый"
		case 6:
			ds = "направление не меняло свое состояние, зеленый"
		case 7:
			ds = "направление не меняло свое состояние, красный"
		case 8:
			ds = "зеленый мигающий сигнал"
		case 9:
			ds = "желтый мигающий в режиме ЖМ"
		case 10:
			ds = "сигналы выключены в режиме ОС"
		case 11:
			ds = "неиспользуемое направление"
		default:
			ds = "error code"
		}
		content = append(content, []any{i, st, ds, hs.TOOBs[i]})
		count++
	}
	rui.SetParams(view, "idNaps", rui.Params{
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}
func updateHeader(view rui.View) {
	t := time.Now()
	rui.Set(view, "time", "text", fmt.Sprintf("<b>%d %02d:%02d:%02d</b>", setup.Set.Id,
		t.Hour(), t.Minute(), t.Second()))
	if hardware.StateHardware.GetConnectUtopia() {
		rui.Set(view, "isCentral", "background-color", "green")
	} else {
		rui.Set(view, "isCentral", "background-color", "gray")
	}
	if hardware.StateHardware.GetConnect() {
		rui.Set(view, "isKDM", "background-color", "green")
	} else {
		rui.Set(view, "isKDM", "background-color", "gray")
	}
	if setup.Set.ModbusRadar.Radar {
		if radar.GetWork() {
			rui.Set(view, "isRadar", "background-color", "green")
		} else {
			rui.Set(view, "isRadar", "background-color", "gray")
		}
	}
	if setup.Set.TrafficData.Work {
		if traffic.GetWork() {
			rui.Set(view, "isTrafficData", "background-color", "green")
		} else {
			rui.Set(view, "isTrafficData", "background-color", "gray")
		}
	}
}
func makeMainScreen(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	updatePartKDM(view)
	updateHeader(view)
	updateMessages(view)

}
func updaterScreen(view rui.View, session rui.Session) {
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
		makeMainScreen(view)
	}
}

func mainScreen(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, mainText)
	if view == nil {
		return nil
	}
	makeMainScreen(view)
	makeButtonOnScreen(view)
	go updaterScreen(view, session)

	return view
}
