package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/potop/hardware"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/stcip/transport"
)

const noSTCIPText = `
ListLayout {
	orientation = vertical, style = showPage,
	content = [
		TextView {
			style=header1,
			text = "STCIP отключена в текущей конфигурации"
		},
	]
}	
`

const stcipText = `
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
							style=table,semantic=code,
							row = 0, column = 0,cell-horizontal-align = left,  title = "Команды", id=idCommands,
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

					]
				},
			]
		},
		TableView {
			style=table,
			cell-horizontal-align = left,  title = "Ответы", id=idReplays,
		},
	]
}
`

func makeScreenSTCIP(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()

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
	rui.Set(view, "idLine1", "text", fmt.Sprintf("<b>Последняя команда в</b> %s",
		toString(hs.LastOperation)))
	rui.Set(view, "idLine2", "text", fmt.Sprintf("<b>План</b> %s <b>Фаза</b> %s <b>Статус КДМ</b> % 02X ",
		PlanToString(hs.Plan), PhaseToString(hs.Phase), hs.Status))
	rui.Set(view, "idLine3", "text", fmt.Sprintf("<b>Расшифровка статуса : %s </b>", hardware.GetError()))

	updateHeader(view)
	var commands [][]any
	commands = append(commands, []any{"Время", "Сообщение"})
	count := 1
	for _, v := range transport.GetHistory(1) {
		commands = append(commands, []any{toString(v.Time), v.Message})
		count++
	}
	rui.SetParams(view, "idCommands", rui.Params{
		rui.Content:             commands,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})
	var reps [][]any
	reps = append(reps, []any{"Время", "Сообщение"})
	count = 1
	for _, v := range transport.GetHistory(2) {
		reps = append(reps, []any{toString(v.Time), v.Message})
		count++
	}
	rui.SetParams(view, "idReplays", rui.Params{
		rui.Content:             reps,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}

func makeButtonOnScreenSTCIP(view rui.View) {
	rui.Set(view, "setAllRedOn", rui.ClickEvent, func(rui.View) {
		transport.CommandFromWeb <- transport.Command{Code: transport.CodeCallAllRed, Value: 1}
		logger.Info.Printf("Оператор установил КК")
	})
	rui.Set(view, "setFlashingOn", rui.ClickEvent, func(rui.View) {
		transport.CommandFromWeb <- transport.Command{Code: transport.CodeCallFlash, Value: 1}
		logger.Info.Printf("Оператор установил ЖМ")
	})
	rui.Set(view, "setDarkOn", rui.ClickEvent, func(rui.View) {
		transport.CommandFromWeb <- transport.Command{Code: transport.CodeCallDark, Value: 1}
		logger.Info.Printf("Оператор установил ОС")
	})
	rui.Set(view, "setPlan", rui.ClickEvent, func(rui.View) {
		transport.CommandFromWeb <- transport.Command{Code: transport.CodeCallPlan, Value: getInteger(rui.Get(view, "idPlan", "value"))}
		logger.Info.Printf("Оператор вызвал план %d", getInteger(rui.Get(view, "idPlan", "value")))
	})
	rui.Set(view, "setPhase", rui.ClickEvent, func(rui.View) {
		transport.CommandFromWeb <- transport.Command{Code: transport.CodeCallPhase, Value: getInteger(rui.Get(view, "idPhase", "value"))}
		logger.Info.Printf("Оператор вызвал фазу %d", getInteger(rui.Get(view, "idPhase", "value")))
	})
}

func updaterSTCIP(view rui.View, session rui.Session) {
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
		makeScreenSTCIP(view)
	}
}

func stcipScreen(session rui.Session) rui.View {
	if setup.Set.STCIP.Run {
		view := rui.CreateViewFromText(session, stcipText)
		if view == nil {
			return nil
		}
		makeScreenSTCIP(view)
		makeButtonOnScreenSTCIP(view)
		go updaterSTCIP(view, session)
		return view
	}
	view := rui.CreateViewFromText(session, noSTCIPText)
	return view

}
