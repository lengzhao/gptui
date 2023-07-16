package window

import (
	"encoding/json"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/lengzhao/gptui/event"
)

type AppConfig struct {
	Theme         string `json:"theme,omitempty"`
	GPTType       string `json:"gptType,omitempty"`
	AzureApiKey   string `json:"azureApiKey,omitempty"`
	AzureEndpoint string `json:"azureEndpoint,omitempty"`
	Model         string `json:"model,omitempty"`
	OpenaiApiKey  string `json:"OPENAI_API_KEY,omitempty"`
	HistoryLimit  int    `json:"HistoryLimit,omitempty"`
	FyneFont      string `json:"FYNE_FONT,omitempty"`
	HttpsProxy    string `json:"HTTPS_PROXY,omitempty"`
	HttpTimeout   int    `json:"HTTP_TIMEOUT,omitempty"`
	PromptsDir    string `json:"PROMPTS_DIR,omitempty"`
	WindWidth     int    `json:"windWidth,omitempty"`
	WindHeight    int    `json:"windHeight,omitempty"`
	Prompt        string `json:"prompt,omitempty"`
}

var confFile string = "config.json"

func makeFormTab() fyne.CanvasObject {
	var conf AppConfig = AppConfig{
		Theme:         "default",
		GPTType:       "openai",
		Model:         "gpt-3.5-turbo",
		AzureEndpoint: "https://openai-poc-instance-east-us.openai.azure.com",
		HistoryLimit:  3,
		HttpTimeout:   20,
		PromptsDir:    "prompts",
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

	prompt := widget.NewEntry()
	prompt.SetText(conf.PromptsDir)
	prompt.OnChanged = func(s string) {
		conf.PromptsDir = s
	}

	label := widget.NewLabel("please restart to take effect")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: themeW, HintText: "App Theme"},
			{Text: "Window Width", Widget: wWidth, HintText: "width of window"},
			{Text: "Window Hight", Widget: wHeight, HintText: "height of window"},
			{Text: "GPT Type", Widget: gptTypeW, HintText: "gpt type"},
			{Text: "Openai Key", Widget: key1, HintText: "the key for openai api"},
			{Text: "Azure Key", Widget: key2, HintText: "the key for azure api"},
			{Text: "Azure Endpoint", Widget: azureEnd, HintText: "the endpoint for azure api"},
			{Text: "Model", Widget: model, HintText: "the model for gpt"},
			{Text: "Chat History Limit", Widget: historyLimit, HintText: "history limit, the history will send to gpt"},
			{Text: "HTTPS Proxy", Widget: httpProxy, HintText: "the https proxy"},
			{Text: "HTTP Timeout", Widget: timeout, HintText: "the http timeout"},
			{Text: "Prompt Dir", Widget: prompt, HintText: "prompt dir, support .json and .csv"},
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
