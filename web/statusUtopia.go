package web

import (
	"time"

	"github.com/anoshenko/rui"
)

const utopiaText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		TextView {
			text = "<b>Текущее состояние Utopia </b>",text-align="center",text-size="24px",
		},
		TextView {
			id=idLine1,text = "Line1",text-size="24px",
		},
		TextView {
			id=idLine2,text = "Line2",text-size="24px",
		},
		TextView {
			id=idLine3,text = "Line3",text-size="24px",
		},
		TextView {
			id=idLine4,text = "Line4",text-size="24px",
		},
		TextView {
			id=idLastOperation,text = "<b>Последняя команда Utopia на устройстве  </b>",text-align="center",text-size="24px",
		},
	]
}
`

func makeViewUtopia(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()

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
	view := rui.CreateViewFromText(session, utopiaText)
	if view == nil {
		return nil
	}
	makeViewUtopia(view)
	go updaterUtopia(view, session)

	return view

}
