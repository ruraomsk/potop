package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/potop/hardware"
)

// TextView {
// 	id=idLine3,text = "Line3",text-size="24px",
// },

const KDMText = `
ListLayout {
	style = showPage,
	orientation = vertical,
	content = [
		TextView {
			style=header1,
			id=idHeader,text = "<b>Текущее состояние КДМ </b>",
		},
		TextView {
			id=idLine1,text = "Line1",text-size="18px",
		},
		TextView {
			id=idLine2,text = "Line2",text-size="18px",
		},
		TableView {
			cell-horizontal-align = left,  title = "Направления", id=idNaps,
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
	source := "внутренний"
	if hs.SourceTOOB {
		source = "внешний"
	}
	rui.Set(view, "idHeader", "text", fmt.Sprintf("<b>Текущее состояние КДМ %s</b>", toString(time.Now())))
	rui.Set(view, "idLine1", "text", fmt.Sprintf("<b>Связь с КДM</b> %v  <b>Utopia</b> %v <b>Последняя команда в</b> %s <b>Тмин=%d Маска=%x остаток watchdog=%d</b>", toRussian(hs.Connect), toRussian(hs.Utopia), toString(hs.LastOperation), hs.Tmin, hs.MaskCommand, hs.RealWatchDog))
	rui.Set(view, "idLine2", "text", fmt.Sprintf("<b>OC</b> %v  <b>KK</b> %v <b>ЖМ</b> %v <b>WatchDog</b> %d <b>План</b> %d <b>Статус КДМ</b> % 02X <b>Источник ТООВ</b> %s",
		toRussian(hs.Dark), toRussian(hs.AllRed), toRussian(hs.Flashing),
		hs.WatchDog, hs.Plan, hs.Status,
		source))
	// rui.Set(view, "idLine3", "text", fmt.Sprintf("<b>Направления % 02X </b>", hs.StatusDirs))
	var content [][]any
	content = append(content, []any{"Нап", "Задание", "Состояние", "Счетчик ТООВ"})
	count := 1
	s := hs.MaskCommand
	for i := 0; i < 32; i++ {
		st := "Закрыто"
		if s&0x80000000 != 0 {
			st = "Открыто"
		}
		s = s << 1
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
