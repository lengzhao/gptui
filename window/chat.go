package window

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/gptui/event"
)

func makeHistoryBox() fyne.CanvasObject {
	var history []string
	dialog := widget.NewMultiLineEntry()
	dialog.Wrapping = fyne.TextWrapWord

	event.RegistEvent(event.EAll, func(key event.EventID, info string) {
		switch key {
		case event.EHistory:
			history = append(history, info)
			if len(history) > 30 {
				history = history[1:]
			}
			dialog.SetText(strings.Join(history, "\n"))
			dialog.CursorRow = len(dialog.Text) - 1
		case event.EUserCommit:
			event.SendEvent(event.EHistory, "---------------------------\nUser: "+info)
		case event.EChatError:
			event.SendEvent(event.EHistory, "GPT Error: "+info)
		case event.EFinishChat:
			event.SendEvent(event.EHistory, "AI: "+info)
		}
	})

	return dialog
}

func makeInputBox() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapBreak
	entry.SetMinRowsVisible(5)
	// entry.SetText("你好,世界!")

	btn := widget.NewButton("Send!", func() {
		if entry.Text == "" {
			return
		}
		event.SendEvent(event.EUserCommit, entry.Text)
		entry.SetText("")
	})

	return container.NewVBox(widget.NewLabel("Input"), entry, btn)
}

func MakeChatWindow() fyne.CanvasObject {
	return container.NewBorder(nil, makeInputBox(), nil, nil, makeHistoryBox())
}
