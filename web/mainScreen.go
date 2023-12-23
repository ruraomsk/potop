package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/diagramm"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/journal"
	"github.com/ruraomsk/potop/radar"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/stcip"
	"github.com/ruraomsk/potop/traffic"
	"github.com/ruraomsk/potop/utopia"
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
					id="isAutonom",text-align = center,
					text = "АВТОНОМ",background-color=gray,
				},
				TextView {
					row = 0, column = 5,
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
							row = 5, column = 0,
							id=setAutoOn,content="Автоном вкл"
						},
						Button {
							row = 5, column = 1,
							id=setAutoOff,content="Автоном выкл"
						},
						Button {
							row = 6, column = 0:1,
							id=setPlan,content="Установить План"
						},
						NumberPicker {
							row = 6, column = 2,
							id=idPlan,type=editor,min=0,max=32,value=0
						},
						Button {
							row = 7, column = 0:1,
							id=setPhase,content="Установить Фазу"
						},
						NumberPicker {
							row = 7, column = 2,
							id=idPhase,type=editor,min=0,max=32,value=0
						},
						Button {
							row = 8, column = 0:2,
							id=setViewDiagramm,content="Показать диаграмму"
						},

					]
				},
			]
		},
		TableView {
			style=table,semantic=code,
			cell-horizontal-align = left,  title = "Диаграма", id=idDiagramm,
		},
		TableView {
			style=table,
			cell-horizontal-align = left,  title = "Диагностика", id=idDiags,
		},
	]
}
`

var viewDiagramm = false

func PlanToString(plan int) string {
	if plan == 0 {
		return "ЛПУ"
	}
	if plan == 25 {
		return "ОС"
	}
	if plan == 26 {
		return "КК"
	}
	if plan == 27 {
		return "ЖМ"
	}
	if plan == 28 {
		return "ВУ"
	}
	return fmt.Sprintf("%d", plan)
}
func PhaseToString(phase int) string {
	if phase == 0 {
		return "KK"
	}
	if phase == 33 {
		return "ЖМ"
	}
	if phase == 34 {
		return "ОС"
	}
	if phase == 35 || phase == 36 {
		return "ВУ"
	}
	if phase == 255 {
		return "ПТ"
	}
	return fmt.Sprintf("%d", phase)

}
func updateDiagramm(view rui.View) {
	var content [][]any
	h := make([]any, 0)
	h = append(h, "Канал")
	var ms []diagramm.ViewDiag
	var times []time.Time
	if viewDiagramm {
		ms, times = diagramm.GetDiagramm(time.Now().Add(-80*time.Second), time.Now())
		for _, v := range times {
			h = append(h, v.Second())
		}
		for i := 0; i < 80-len(times); i++ {
			h = append(h, "  ")
		}
	}

	content = append(content, h)
	count := 1
	if viewDiagramm {
		for _, v := range ms {
			if v.View {
				r := make([]any, 0)
				r = append(r, v.Chanel)
				for _, k := range v.Data {
					r = append(r, k)
				}
				for i := 0; i < 80-len(r); i++ {
					r = append(r, "__")
				}
				content = append(content, r)
				count++
			}
		}

	}
	rui.SetParams(view, "idDiagramm", rui.Params{
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}
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
		hardware.CommandToKDM(4, 1)
		logger.Info.Printf("Оператор установил КК")
	})
	rui.Set(view, "setFlashingOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(3, 1)
		logger.Info.Printf("Оператор установил ЖМ")
	})
	rui.Set(view, "setDarkOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(6, 1)
		logger.Info.Printf("Оператор установил ОС")
	})
	rui.Set(view, "setAutoOn", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(0, 1)
		if setup.Set.Utopia.Run {
			utopia.SetAutonom(true)
		}
		if setup.Set.STCIP.Run {
			stcip.SetAutonom(true)
		}
		logger.Info.Printf("Оператор перевел в Автоном")
	})
	rui.Set(view, "setAutoOff", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(0, 0)
		if setup.Set.Utopia.Run {
			utopia.SetAutonom(false)
		}
		if setup.Set.STCIP.Run {
			stcip.SetAutonom(false)
		}
		logger.Info.Printf("Оператор отключил Автоном")
	})
	rui.Set(view, "setPlan", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(7, getInteger(rui.Get(view, "idPlan", "value")))
		logger.Info.Printf("Оператор вызвал план %d", getInteger(rui.Get(view, "idPlan", "value")))
	})
	rui.Set(view, "setPhase", rui.ClickEvent, func(rui.View) {
		hardware.CommandToKDM(8, getInteger(rui.Get(view, "idPhase", "value")))
		logger.Info.Printf("Оператор вызвал фазу %d", getInteger(rui.Get(view, "idPhase", "value")))
	})
	rui.Set(view, "setViewDiagramm", rui.ClickEvent, func(rui.View) {
		if viewDiagramm {
			viewDiagramm = false
			rui.Set(view, "setViewDiagramm", "content", "Показать диаграмму")
		} else {
			viewDiagramm = true
			rui.Set(view, "setViewDiagramm", "content", "Скрыть диаграмму")
		}
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
	rui.Set(view, "idLine1", "text", fmt.Sprintf("<b>Последняя команда в</b> %s <b>Тмин=%d Маска=%x остаток watchdog=%d</b>",
		toString(hs.LastOperation), hs.Tmin, hs.MaskCommand, hs.RealWatchDog))
	rui.Set(view, "idLine2", "text", fmt.Sprintf("<b>WatchDog</b> %d <b>План</b> %s <b>Фаза</b> %s <b>Статус КДМ</b> % 02X <b>Источник ТООВ</b> %s",
		hs.WatchDog, PlanToString(hs.Plan), PhaseToString(hs.Phase), hs.Status, source))
	rui.Set(view, "idLine3", "text", fmt.Sprintf("<b>Расшифровка статуса : %s </b>", hardware.GetError()))
	var content [][]any
	if hs.SourceTOOB {
		content = append(content, []any{"Нап", "Задание", "Состояние", "Счетчик ТООВ"})
	} else {
		content = append(content, []any{"Нап", "Состояние"})
	}
	count := 1
	s := hs.MaskCommand
	for i := 0; i < 32; i++ {
		st := "Закрыто"
		if s&0x1 != 0 {
			st = "Открыто"
		}
		s = s >> 1
		ds := "undef"
		if hs.StatusDirs[i] == 11 {
			continue
		}
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
		if hs.SourceTOOB {
			content = append(content, []any{i, st, ds, hs.TOOBs[i]})
		} else {
			content = append(content, []any{i, ds})
		}
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
	if hardware.StateHardware.GetCentral() {
		rui.Set(view, "isCentral", "background-color", "green")
	} else {
		rui.Set(view, "isCentral", "background-color", "gray")
	}
	if hardware.StateHardware.GetConnect() {
		rui.Set(view, "isKDM", "background-color", "green")
	} else {
		rui.Set(view, "isKDM", "background-color", "gray")
	}
	if setup.Set.ModbusRadar.Work {
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
	if setup.Set.Utopia.Run {
		if utopia.GetAutonom() {
			rui.Set(view, "isAutonom", "background-color", "green")
		} else {
			rui.Set(view, "isAutonom", "background-color", "gray")
		}
	}
	if setup.Set.STCIP.Run {
		if stcip.GetAutonom() {
			rui.Set(view, "isAutonom", "background-color", "green")
		} else {
			rui.Set(view, "isAutonom", "background-color", "gray")
		}
	}
}
func makeMainScreen(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	updatePartKDM(view)
	updateHeader(view)
	updateMessages(view)
	updateDiagramm(view)
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
