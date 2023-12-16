package web

import (
	"fmt"

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
			id=idHeader,text = "<b>Текущие привязки КДМ </b>",
		},
		TabsLayout { id = tabsPK, width = 100%, height = 100%, tabs = top, tab-close-button = false,
			content = [
				TableView {cell-horizontal-align = right,  title = "Конфликты", id="idConflicts"},
				TableView {cell-horizontal-align = right,  title = "Направления", id="idNaps"},
				TableView {cell-horizontal-align = right,  title = "Фазы", id="idPhases"},
				TableView {cell-horizontal-align = right,  title = "Планы", id="idPlans"},
				TableView {cell-horizontal-align = right,  title = "Годовой", id="idYears"},
				TableView {cell-horizontal-align = right,  title = "Недельный", id="idWeeks"},
				TableView {cell-horizontal-align = right,  title = "Суточный", id="idDays"},
			]
		}
	]
}
`

var config hardware.Config

func typePhase(t int) string {
	switch t {
	case 0:
		return "простая"
	case 1:
		return "МГР"
	case 2:
		return "ТВП1"
	case 3:
		return "ТВП2"
	case 4:
		return "ТВП 1,2"
	case 5:
		return "ЗАМ 1"
	case 6:
		return "ЗАМ 2"
	case 7:
		return "ЗАМ 1,2"
	}
	return fmt.Sprintf("%d", t)
}
func typeNap(t int) string {
	switch t {
	case 0:
		return "Не назначен"
	case 1:
		return "Транспортное"
	case 2:
		return "Пешеходное"
	case 3:
		return "Стрелка"
	}
	return fmt.Sprintf("%d", t)
}
func makeConflicts(view rui.View) {
	var content [][]any
	ph := []any{"Номер"}
	for i := 1; i < 33; i++ {
		ph = append(ph, i)
	}
	content = append(content, ph)
	for i := 0; i < 31; i++ {
		lh := []any{i + 1}

		for _, v := range config.Conflicts.Conflicts[i].Line {
			if v {
				lh = append(lh, 1)
			} else {
				lh = append(lh, 0)
			}
		}
		content = append(content, lh)
	}
	rui.SetParams(view, "idConflicts", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}
func makeNaps(view rui.View) {
	var content [][]any
	ph := []any{"Номер", "Тип", "Змиг", "Жп", "Кп", "КЖ", "Табло"}
	for i := 1; i < 65; i++ {
		ph = append(ph, i)
	}
	content = append(content, ph)
	for i := 0; i < 31; i++ {
		lh := []any{i + 1, typeNap(config.DefineNaps.DefineNaps[i].Type),
			config.DefineNaps.DefineNaps[i].FlashGreen,
			config.DefineNaps.DefineNaps[i].YellowFlash,
			config.DefineNaps.DefineNaps[i].Red,
			config.DefineNaps.DefineNaps[i].RedYellow,
			config.DefineNaps.DefineNaps[i].Tiob}
		for _, v := range config.DefineNaps.DefineNaps[i].Keys {
			if v {
				lh = append(lh, 1)
			} else {
				lh = append(lh, 0)
			}
		}
		content = append(content, lh)
	}
	rui.SetParams(view, "idNaps", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})

}
func makePhases(view rui.View) {
	var content [][]any
	ph := []any{"Номер фазы", "Тмин"}
	for i := 1; i < 33; i++ {
		ph = append(ph, i)
	}
	content = append(content, ph)
	for i := 0; i < 31; i++ {
		lh := []any{i + 1, config.DefPhases.DefPhases[i].Tmin}
		for _, v := range config.DefPhases.DefPhases[i].Naps {
			if v {
				lh = append(lh, 1)
			} else {
				lh = append(lh, 0)
			}
		}
		content = append(content, lh)
	}
	rui.SetParams(view, "idPhases", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "5px",
		rui.CellHorizontalAlign: "left",
	})

}
func makePlans(view rui.View) {
	var content [][]any
	content = append(content, []any{"Номер плана", "Вид", "Сдвиг", "Дительность", "Фаза", "Тип", "Начало", "Конец"})
	//RPU
	plan := config.RPU
	v := "ЛПУ"
	if plan.Type == 1 {
		v = "КРУ"
	}
	content = append(content, []any{"РПУ", v, plan.Shift, plan.Tcycle,
		plan.Lines[0].Phase, typePhase(plan.Lines[0].Type), plan.Lines[0].Start, plan.Lines[0].Stop})
	for j := 1; j < 24; j++ {
		if (plan.Lines[j].Phase == 0) && (plan.Lines[j].Start == 0) && (plan.Lines[j].Stop == 0) {
			break
		}
		content = append(content, []any{"", "", "", "", plan.Lines[j].Phase, typePhase(plan.Lines[j].Type), plan.Lines[j].Start, plan.Lines[j].Stop})
	}
	content = append(content, []any{"", "", "", "", "", "", "", ""})

	for i := 0; i < 24; i++ {
		plan := config.Plans.Plans[i]
		if len(plan.Lines) == 0 {
			continue
		}
		if (plan.Tcycle == 0) && (plan.Shift == 0) && (plan.Type == 0) {
			continue
		}
		v := "ЛПУ"
		if plan.Type == 1 {
			v = "КРУ"
		}
		content = append(content, []any{i + 1, v, plan.Shift, plan.Tcycle,
			plan.Lines[0].Phase, typePhase(plan.Lines[0].Type), plan.Lines[0].Start, plan.Lines[0].Stop})
		for j := 1; j < 24; j++ {
			if (plan.Lines[j].Phase == 0) && (plan.Lines[j].Start == 0) && (plan.Lines[j].Stop == 0) {
				break
			}
			content = append(content, []any{"", "", "", "", plan.Lines[j].Phase, typePhase(plan.Lines[j].Type), plan.Lines[j].Start, plan.Lines[j].Stop})
		}
	}
	rui.SetParams(view, "idPlans", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "1px",
		rui.CellHorizontalAlign: "left",
	})
}

func makeDays(view rui.View) {
	var content [][]any
	content = append(content, []any{"Номер карты", "Старт", "Финиш", "План"})
	for i := 0; i < 12; i++ {
		l := config.Days.Days[i].Nplans
		if len(l) == 0 {
			continue
		}
		if (l[0].Plan == 0) && (l[0].Start == 0) && (l[0].Stop == 0) {
			continue
		}
		content = append(content, []any{i + 1, l[0].Start, l[0].Stop, l[0].Plan})
		for j := 1; j < 12; j++ {
			if (l[j].Plan == 0) && (l[j].Start == 0) && (l[j].Stop == 0) {
				break
			}
			content = append(content, []any{"", l[j].Start, l[j].Stop, l[j].Plan})
		}
	}
	rui.SetParams(view, "idDays", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "10px",
		rui.CellHorizontalAlign: "left",
	})
}

func makeWeeks(view rui.View) {
	var content [][]any
	content = append(content, []any{"Номер", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"})
	for i := 0; i < 12; i++ {
		l := config.Weeks.Weeks[i]
		content = append(content, []any{i + 1, l.Days[0], l.Days[1], l.Days[2], l.Days[3], l.Days[4], l.Days[5], l.Days[6]})
	}
	rui.SetParams(view, "idWeeks", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "10px",
		rui.CellHorizontalAlign: "left",
	})
}
func makeYears(view rui.View) {
	var content [][]any
	header := []any{"Месяц"}
	for i := 1; i < 32; i++ {
		header = append(header, i)
	}
	content = append(content, header)
	for i := 0; i < 12; i++ {
		l := []any{i + 1}
		for _, v := range config.Year.Year[i].Days {
			l = append(l, v)
		}
		content = append(content, l)
	}
	rui.SetParams(view, "idYears", rui.Params{
		rui.Style:               "table",
		rui.Content:             content,
		rui.HeadHeight:          1,
		rui.CellPadding:         "5px",
		rui.CellHorizontalAlign: "left",
	})
}

func makeViewKDM(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	config = hardware.GetConfig()
	if !config.Ready {
		return
	}
	makeWeeks(view)
	makeYears(view)
	makeDays(view)
	makePlans(view)
	makePhases(view)
	makeNaps(view)
	makeConflicts(view)
}

func statusKDM(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, KDMText)
	if view == nil {
		return nil
	}
	makeViewKDM(view)
	return view

}
