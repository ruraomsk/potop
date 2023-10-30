package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/potop/hardware"
)

const KDMText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		TextView {
			id=idHeader,text = "<b>Текущее состояние КДМ </b>",text-align="center",text-size="24px",
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
	]
}
`

func toRussian(t bool) string {
	if t {
		return "есть"
	} else {
		return "нет"
	}

}
func makeViewKDM(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	hs := hardware.GetStateHard()
	rui.Set(view, "idHeader", "text", fmt.Sprintf("<b>Текущее состояние КДМ %s</b>", toString(time.Now())))
	rui.Set(view, "idLine1", "text", fmt.Sprintf("<b>Связь с КДM</b> %v  <b>Utopia</b> %v <b>Последняя команда в</b> %s", toRussian(hs.Connect), toRussian(hs.Utopia), toString(hs.LastOperation)))
	rui.Set(view, "idLine2", "text", fmt.Sprintf("<b>OC</b> %v  <b>KK</b> %v <b>ЖМ</b> %v <b>WatchDog</b> %d <b>План</b> %d <b>Статус КДМ</b> % 02X ", toRussian(hs.Dark), toRussian(hs.AllRed), toRussian(hs.Flashing), hs.WatchDog, hs.Plan, hs.Status))
	rui.Set(view, "idLine3", "text", fmt.Sprintf("<b>Направления % 02X </b>", hs.StatusDirs))
	rui.Set(view, "idLine4", "text", fmt.Sprintf("<b>Тмин=%d Маска=%x остаток watchdog=%d</b>", hs.Tmin, hs.MaskCommand, hs.RealWatchDog))
}
func updaterKDM(view rui.View, session rui.Session) {
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
		makeViewKDM(view)
	}
}

func statusKDM(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, KDMText)
	if view == nil {
		return nil
	}
	makeViewKDM(view)
	go updaterKDM(view, session)

	return view

}
