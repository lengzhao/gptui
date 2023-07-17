package window

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/conf"
	"github.com/lengzhao/gptui/event"
)

//go:embed embed
var prompts embed.FS

type PromptItem struct {
	Act    string `json:"act,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

func makeListTab() fyne.CanvasObject {
	data := loadPrompts("embed/default.json")

	label := widget.NewMultiLineEntry()
	label.Wrapping = fyne.TextWrapWord
	label.SetPlaceHolder("You can change and set whatever you want.")

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
			fmt.Println("not found prompts,", info)
			return
		}

		data = prompts
		list.Refresh()
		list.Select(0)
		label.SetText(data[0].Prompt)
	})
	btn := widget.NewButton("Set", func() {
		event.SendEvent(event.ESystemPrompt, strings.TrimSpace(label.Text))
	})
	out := container.NewHSplit(list, container.NewBorder(nil, btn, nil, nil, label))

	return out
}

func makeSourceTab(win fyne.Window) fyne.CanvasObject {
	var promptFiles = []string{"embed/default.json", "embed/prompts-en.csv", "embed/prompts-zh.json"}
	dir := conf.Get("PROMPTS_DIR", "prompts")
	os.Mkdir(dir, 0766)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".csv" && ext != ".json" {
			return nil
		}

		promptFiles = append(promptFiles, path)
		return nil
	})

	selectEntry := widget.NewSelect(promptFiles, func(s string) {
		event.SendEvent(event.ELoadPrompt, s)
	})
	if len(promptFiles) > 0 {
		selectEntry.SetSelected(promptFiles[0])
	}
	openFile := widget.NewButton("Load Customized Prompt(.json/.csv)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			fmt.Println("path", reader.URI().Path())
			event.SendEvent(event.ELoadPrompt, reader.URI().Path())
		}, win)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json", ".csv"}))
		fd.Show()
	})

	return container.NewBorder(nil, nil, nil, openFile, selectEntry)
}

func loadPrompts(fn string) []PromptItem {
	f, err := os.Open(fn)
	if err != nil {
		// fmt.Println("fail to open file:", fn, err)
		return loadFromEmbed(fn)
	}
	defer f.Close()
	return loadFromFile(f, filepath.Ext(fn))
}

func MakeRoleWindow(win fyne.Window) fyne.CanvasObject {
	return container.NewBorder(makeSourceTab(win), nil, nil, nil, makeListTab())
}

func loadFromEmbed(fn string) []PromptItem {
	f, err := prompts.Open(fn)
	if err != nil {
		return nil
	}
	defer f.Close()

	return loadFromFile(f, filepath.Ext(fn))
}

func loadFromFile(f io.Reader, ext string) []PromptItem {
	var out []PromptItem
	switch ext {
	case ".json":
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&out)
		if err != nil {
			fmt.Println("fail to Decode file:", err)
			return nil
		}
	case ".csv":
		items, err := csv.NewReader(f).ReadAll()
		if err != nil {
			fmt.Println("fail to Read file:", err)
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
		fmt.Println("unknown file type:", ext)

	}
	return out
}
