package window

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/conf"
	"github.com/lengzhao/gptui/event"
)

type PromptItem struct {
	Act    string `json:"act,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

func makeListTab() fyne.CanvasObject {
	data := loadPrompts("default")

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
			item.(*widget.Label).SetText(data[id].Act)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id].Prompt)
	}
	list.Select(0)
	event.RegistEvent(event.ELoadPrompt, func(key event.EventID, info string) {
		prompts := loadPrompts(info)
		if len(prompts) == 0 {
			return
		}

		data = prompts
		list.Refresh()
		list.Select(0)
		label.SetText(data[0].Prompt)
	})
	btn := widget.NewButton("Set Prompt", func() {
		event.SendEvent(event.ESystemPrompt, label.Text)
	})
	out := container.NewHSplit(list, container.NewBorder(nil, btn, nil, nil, label))

	return out
}

func makeSourceTab() fyne.CanvasObject {
	var files []string = []string{"default"}
	dir := conf.Get("PROMPTS_DIR", "prompts")
	os.Mkdir(dir, 0666)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".csv" && ext != ".json" {
			return nil
		}

		files = append(files, path)
		return nil
	})
	selectEntry := widget.NewSelect(files, func(s string) {
		// fmt.Println("change prompt file:", s)
		event.SendEvent(event.ELoadPrompt, s)
	})
	if len(files) > 0 {
		selectEntry.SetSelected(files[0])
	}

	return selectEntry
}

func loadPrompts(fn string) []PromptItem {
	var out []PromptItem
	if fn == "default" {
		out = append(out, PromptItem{
			Act:    "AI Writing Tutor",
			Prompt: `I want you to act as an AI writing tutor. I will provide you with a student who needs help improving their writing and your task is to use artificial intelligence tools, such as natural language processing, to give the student feedback on how they can improve their composition. You should also use your rhetorical knowledge and experience about effective writing techniques in order to suggest ways that the student can better express their thoughts and ideas in written form. My first request is "I need somebody to help me edit my master's thesis."`,
		})
		return out
	}
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println("fail to open file:", fn, err)
		return nil
	}
	defer f.Close()
	switch filepath.Ext(fn) {
	case ".json":
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&out)
		if err != nil {
			fmt.Println("fail to Decode file:", fn, err)
			return nil
		}
	case ".csv":
		items, err := csv.NewReader(f).ReadAll()
		if err != nil {
			fmt.Println("fail to Read file:", fn, err)
			return nil
		}
		for _, item := range items {
			if len(item) < 2 {
				continue
			}
			if item[0] == "act" && item[1] == "prompt" {
				continue
			}
			out = append(out, PromptItem{Act: item[0], Prompt: item[1]})
		}
	default:
		fmt.Println("unknown file type:", fn, filepath.Ext(fn))

	}
	return out
}

func MakeRoleWindow() fyne.CanvasObject {
	return container.NewBorder(makeSourceTab(), nil, nil, nil, makeListTab())
}
