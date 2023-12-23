package web

import (
	"fmt"
	"time"

	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/utopia"

	"github.com/anoshenko/rui"
)

// row=1,column=0,
// width = 100%, height = 100%, cell-width = "150px, 1fr, 30%", cell-height = "25%, 200px, 1fr",
const noUtopiaText = `
ListLayout {
	orientation = vertical, style = showPage,
	content = [
		TextView {
			style=header1,
			text = "Utopia отключена в текущей конфигурации"
		},
	]
}	
`

const utopiaText = `
ListLayout {
		 orientation = vertical,
		text-align = left, vertical-align = top,
		style = showPage,
		content = [
			TextView {
				style=header1,
				id=idHeader, text = "<b>Текущее состояние Utopia </b>",
				},
				TextView {
					text = "Команды СПОТ",text-size="24px",
					text-align = center,
				},
				TextView {
					id=idLine1,text = "",text-size="16px",
				},
				TextView {
					id=idLine2,text = "",text-size="16px",
				},
				TextView {
					id=idLine3,text = "",text-size="16px",
				},
				TextView {
					id=idLine4,text = "",text-size="16px",
				},
				TextView {
					text = "Ответы контроллера",text-size="24px",
					text-align = center,						
					
				},
				TextView {
					id=idLine11,text = "",text-size="16px",
				},
				TextView {
					id=idLine12,text = "",text-size="16px",
				},
				TextView {
					id=idLine13,text = "",text-size="16px",
				},
				TextView {
					id=idLine14,text = "",text-size="16px",
				},
		]
}
`

func makeViewUtopia(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	ctrl := utopia.GetControllerUtopia()
	rui.Set(view, "idHeader", "text", fmt.Sprintf("<b>Текущее состояние Utopia %s</b>", toString(time.Now())))
	rui.Set(view, "idLine1", "text", ctrl.TlcAndGroupControl.ToString())
	rui.Set(view, "idLine2", "text", ctrl.CountDown.ToString())
	rui.Set(view, "idLine3", "text", ctrl.ExtendedCountDown.ToString())
	rui.Set(view, "idLine11", "text", ctrl.StatusAndDetections.ToString())
	rui.Set(view, "idLine12", "text", ctrl.SignalGroupFeedback.ToString())

}

func updaterUtopia(view rui.View, session rui.Session) {
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
		makeViewUtopia(view)
	}
}

func statusUtopia(session rui.Session) rui.View {
	if setup.Set.Utopia.Run {
		view := rui.CreateViewFromText(session, utopiaText)
		if view == nil {
			return nil
		}
		makeViewUtopia(view)
		go updaterUtopia(view, session)
		return view
	}
	view := rui.CreateViewFromText(session, noUtopiaText)
	return view
}
