package web

import (
	"fmt"
	"strconv"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/potop/setup"
	"github.com/ruraomsk/potop/traffic"
)

const noTrafficText = `
ListLayout {
	orientation = vertical, style = showPage,
	content = [
		TextView {
			style=header1,
			text = "TrafficData отключена в текущей конфигурации"
		},
	]
}	
`
const trafficText = `
		ListLayout {
			width = 100%, height = 100%, orientation = vertical, padding = 16px,
			content = [
				TextView {
					style=header1,
					id=titleTraffic,text = ""
				},
				TextView {
					id=idTimeTraffic,text-size="24px",
					text = "Время на сервере и время upTime"
				},
				TextView {
					id=idConnectionInformation,text-size="24px",
					text = "ConnectionInformation"
				},
				ListLayout {
					width = 100%, height = 100%, orientation = vertical, padding = 16px,
					content = [
						TextView {
							text-align="center",text-size="24px",
							border = _{ style = solid, width = 1px},
							id=titleProduct,text = "Продукт"
						},
						TableView {cell-horizontal-align = right,
							id=tableCameras}
					]
				},
		
			]
		}
`

func makeViewTraffic(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	st := traffic.GetStatusTrafficData()
	t := time.Now()
	if st.StatusWorkClient && st.StatusWorkReciever {
		rui.Set(view, "titleTraffic", "text", fmt.Sprintf("<b>Последний обмен с TrafficData </b>%s",
			st.LastOperation.Format("15:04:05")))
		rui.Set(view, "titleTraffic", "text-color", "green")
	} else {
		rui.Set(view, "titleTraffic", "text", fmt.Sprintf("<b>Отсутствует обмен TrafficData </b>%02d:%02d:%02d",
			t.Hour(), t.Minute(), t.Second()))
		rui.Set(view, "titleTraffic", "text-color", "red")
	}
	ts, _ := strconv.ParseUint(st.StatusGetTime.Body.Utc, 10, 64)

	rui.Set(view, "idTimeTraffic", "text", fmt.Sprintf("<b>Время на сервере</b> %s. <b>Сервер в работе</b> %s секунд",
		time.Unix(int64(ts), 0).Format("15:04:05"), st.StatusUpTime.Body.Value))
	rui.Set(view, "idConnectionInformation", "text", fmt.Sprintf("<b>Server IP</b> %s <b>Subnet</b> %s <b>Gatesay</b> %s",
		st.StatusConnectionInformation.Body.Ip, st.StatusConnectionInformation.Body.Subnet, st.StatusConnectionInformation.Body.Gateway))
	rui.Set(view, "titleProduct", "text", fmt.Sprintf("<b>Компания</b> %s <b>Продукт</b> %s <b>Версия</b> %s",
		st.StatusProductInformation.Body.Company, st.StatusProductInformation.Body.Product, st.StatusProductInformation.Body.Version))

	var content [][]any
	count := 1
	content = append(content, []any{"Камера", "Наименование", "Url", "Streaming"})
	for _, v := range st.StatusProductInformation.Body.Camera {
		count++
		content = append(content, []any{v.Id, v.Name, v.Url, v.StreamingUrl})
	}

	rui.SetParams(view, "tableCameras", rui.Params{
		rui.Content:     content,
		rui.HeadHeight:  count,
		rui.CellPadding: "4px",
	})

}
func updaterTraffic(view rui.View, session rui.Session) {
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
		makeViewTraffic(view)
	}
}

func trafficShow(session rui.Session) rui.View {
	if !setup.Set.TrafficData.Work {
		return rui.CreateViewFromText(session, noTrafficText)
	}
	view := rui.CreateViewFromText(session, trafficText)
	if view == nil {
		return nil
	}
	makeViewTraffic(view)
	go updaterTraffic(view, session)

	return view
}
