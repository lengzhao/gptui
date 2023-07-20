package window

import (
	"encoding/json"
	"os"
	"strconv"

	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/gptui/event"
)

type AppConfig struct {
	Theme         string `json:"theme,omitempty"`
	GPTType       string `json:"gptType"`
	AzureApiKey   string `json:"azureApiKey"`
	AzureEndpoint string `json:"azureEndpoint"`
	Model         string `json:"model,omitempty"`
	OpenaiApiKey  string `json:"OPENAI_API_KEY"`
	HistoryLimit  int    `json:"HistoryLimit,omitempty"`
	FyneFont      string `json:"FYNE_FONT,omitempty"`
	HttpsProxy    string `json:"HTTPS_PROXY,omitempty"`
	HttpTimeout   int    `json:"HTTP_TIMEOUT,omitempty"`
	WindWidth     int    `json:"windWidth,omitempty"`
	WindHeight    int    `json:"windHeight,omitempty"`
	Prompt        string `json:"prompt,omitempty"`
}

//go:generate sh -c "printf %s-%s $(git symbolic-ref HEAD | cut -b 12-) $(git describe --tags --abbrev=8 --dirty --always --long)> version.txt"

//go:embed version.txt
var version string

func makeFormTab() fyne.CanvasObject {
	var confFile string = "config.json"
	var conf AppConfig = AppConfig{
		Theme:         "default",
		GPTType:       "openai",
		Model:         "gpt-3.5-turbo",
		AzureEndpoint: "https://openai-poc-instance-east-us.openai.azure.com",
		HistoryLimit:  3,
		HttpTimeout:   20,
		WindWidth:     1000,
		WindHeight:    800,
	}
	data, err := os.ReadFile(confFile)
	if err == nil {
		json.Unmarshal(data, &conf)
	}
	themeW := widget.NewRadioGroup([]string{"dark", "light", "default"}, func(t string) {
		conf.Theme = t
	})
	themeW.Horizontal = true
	themeW.SetSelected(conf.Theme)

	wWidth := widget.NewEntry()
	wWidth.SetText(strconv.Itoa(conf.WindWidth))
	wWidth.OnChanged = func(s string) {
		conf.WindWidth, _ = strconv.Atoi(s)
	}
	wHeight := widget.NewEntry()
	wHeight.SetText(strconv.Itoa(conf.WindHeight))
	wHeight.OnChanged = func(s string) {
		conf.WindHeight, _ = strconv.Atoi(s)
	}

	gptTypeW := widget.NewRadioGroup([]string{"openai", "azure"}, func(t string) {
		conf.GPTType = t
	})
	gptTypeW.Horizontal = true
	gptTypeW.SetSelected(conf.GPTType)

	key1 := widget.NewPasswordEntry()
	key1.SetPlaceHolder("OpenaiApiKey")
	key1.SetText(conf.OpenaiApiKey)
	key1.OnChanged = func(s string) {
		conf.OpenaiApiKey = s
	}

	key2 := widget.NewPasswordEntry()
	key2.SetPlaceHolder("AzureApiKey")
	key2.SetText(conf.AzureApiKey)
	key2.OnChanged = func(s string) {
		conf.AzureApiKey = s
	}

	azureEnd := widget.NewEntry()
	azureEnd.SetText(conf.AzureEndpoint)
	azureEnd.OnChanged = func(s string) {
		conf.AzureEndpoint = s
	}

	model := widget.NewSelectEntry([]string{"gpt-4",
		"gpt-4-32k", "gpt-3.5-turbo", "gpt-3.5-turbo-16k",
		"text-davinci-003", "text-davinci-002", "code-davinci-002"})
	model.OnChanged = func(s string) {
		conf.Model = s
	}
	model.SetText(conf.Model)

	historyLimit := widget.NewSelect([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, func(s string) {
		conf.HistoryLimit, _ = strconv.Atoi(s)
	})
	historyLimit.SetSelected(strconv.Itoa(conf.HistoryLimit))

	httpProxy := widget.NewEntry()
	httpProxy.SetText(conf.HttpsProxy)
	httpProxy.OnChanged = func(s string) {
		conf.HttpsProxy = s
	}

	timeout := widget.NewSelect([]string{"5", "10", "15", "20", "30", "60"}, func(s string) {
		conf.HttpTimeout, _ = strconv.Atoi(s)
	})
	timeout.SetSelected(strconv.Itoa(conf.HttpTimeout))

	prompt := widget.NewMultiLineEntry()
	prompt.Wrapping = fyne.TextWrapWord
	prompt.SetText(conf.Prompt)
	prompt.OnChanged = func(s string) {
		conf.Prompt = s
	}

	label := widget.NewLabel("please restart to take effect")
	label.TextStyle = fyne.TextStyle{Bold: true, Italic: true, Symbol: true}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Version", Widget: widget.NewLabel(version)},
			{Text: "Theme", Widget: themeW},
			{Text: "Window Width", Widget: wWidth},
			{Text: "Window Hight", Widget: wHeight},
			{Text: "GPT Type", Widget: gptTypeW},
			{Text: "Openai Key", Widget: key1},
			{Text: "Azure Key", Widget: key2},
			{Text: "Azure Endpoint", Widget: azureEnd, HintText: "the endpoint for azure api"},
			{Text: "Model", Widget: model, HintText: "the model for gpt"},
			{Text: "Chat History Limit", Widget: historyLimit, HintText: "history limit, the history will send to gpt"},
			{Text: "HTTPS Proxy", Widget: httpProxy, HintText: "the https proxy"},
			{Text: "HTTP Timeout", Widget: timeout, HintText: "the http timeout"},
			{Text: "Prompt", Widget: prompt, HintText: "prompt for gpt"},
			{Text: "Note", Widget: label},
		},
		OnSubmit: func() {
			// fmt.Println("Form submitted")
			data, err := json.MarshalIndent(conf, "", " ")
			if err != nil {
				return
			}
			os.WriteFile(confFile, data, 0644)
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Need restart app",
				Content: "Restart App for the configuration to take effect",
			})
		},
	}
	event.RegistEvent(event.ESystemPrompt, func(key event.EventID, info string) {
		conf.Prompt = info
		data, err := json.MarshalIndent(conf, "", " ")
		if err != nil {
			return
		}
		os.WriteFile(confFile, data, 0644)
	})
	return form
}

func MakeSettingWindow() fyne.CanvasObject {
	return makeFormTab()
}
