package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/gptui/event"
	"github.com/lengzhao/gptui/history"
)

var data []history.HistoryItem

func init() {
	event.RegistEvent(event.EAll, func(key event.EventID, info string) {
		data = append(data, history.HistoryItem{
			Role: key.String(),
			Text: info,
		})
		if len(data) > 1000 {
			data = data[1:]
		}
	})
}

func MakeLogListTab() fyne.CanvasObject {
	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapWord

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Object")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(data[id].Role)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id].Text)
	}
	list.Select(0)
	btn := widget.NewButton("clear", func() {
		data = []history.HistoryItem{}
	})

	event.RegistEvent(event.EAll, func(key event.EventID, info string) {
		list.Select(len(data) - 1)
		list.Refresh()
	})

	out := container.NewHSplit(list, container.NewBorder(nil, btn, nil, nil, label))

	return out
}
