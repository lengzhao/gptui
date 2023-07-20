package window

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/gptui/chat"
	"github.com/lengzhao/gptui/event"
)

func makeHistoryBox() fyne.CanvasObject {
	var history []string
	dialog := widget.NewMultiLineEntry()
	dialog.Wrapping = fyne.TextWrapWord
	dialog.SetPlaceHolder("History")

	addHist := func(info string) {
		history = append(history, info)
		if len(history) > 30 {
			history = history[1:]
		}
		dialog.SetText(strings.Join(history, "\n"))
		dialog.CursorRow = len(dialog.Text) - 1
	}

	event.RegistEvent(event.EAll, func(key event.EventID, info string) {
		switch key {
		case event.EUserCommit:
			addHist("---------------------------\nUser: " + info)
		case event.EChatError:
			addHist("GPT Error: " + info)
		case event.EFinishChat:
			addHist("AI: " + info)
		case event.EHistory:
			addHist(info)
		}
	})

	return dialog
}

func makeInputBox() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapBreak
	entry.SetMinRowsVisible(5)

	btn := widget.NewButton("Send!", func() {
		if entry.Text == "" {
			return
		}
		event.SendEvent(event.EUserCommit, entry.Text)
		entry.SetText("")
	})

	err := chat.StartWithEvent()
	if err != nil {
		event.SendEvent(event.EError, fmt.Sprintf("fail to start client:%s", err))
	}

	return container.NewVBox(widget.NewLabel("Input"), entry, btn)
}

func MakeChatWindow() fyne.CanvasObject {
	return container.NewBorder(nil, makeInputBox(), nil, nil, makeHistoryBox())
}
