package web

import "github.com/anoshenko/rui"

const controlText = `
		ListLayout {
			style = showPage,
			orientation = vertical,
			content = [
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
					]
				},
			]
		}
`

func controlKDM(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, controlText)
	if view == nil {
		return nil
	}
	return view
}
