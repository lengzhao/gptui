package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/lengzhao/conf"
	_ "github.com/lengzhao/font/autoload"
	"github.com/lengzhao/gptui/window"
)

func main() {
	a := app.New()
	w := a.NewWindow("gpt")
	fn := conf.LoadOneToENV()
	if fn == "" {
		a.SendNotification(&fyne.Notification{
			Title:   "not found config file",
			Content: "not found config file",
		})
	}
	w.Resize(fyne.NewSize(float32(conf.GetFloat("windWidth", 1000)),
		float32(conf.GetFloat("windHeight", 700))))

	w.SetContent(makeTabs(w))
	switch conf.Get("theme", "default") {
	case "dark":
		a.Settings().SetTheme(theme.DarkTheme())
	case "light":
		a.Settings().SetTheme(theme.LightTheme())
	default:
		a.Settings().SetTheme(theme.DefaultTheme())
	}

	w.ShowAndRun()
}

func makeTabs(win fyne.Window) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Chat", theme.AccountIcon(), window.MakeChatWindow()),
		container.NewTabItemWithIcon("Role", theme.DocumentPrintIcon(), window.MakeRoleWindow(win)),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), window.MakeSettingWindow()),
		container.NewTabItemWithIcon("Logger", theme.ContentClearIcon(), window.MakeLogListTab()),
	)
	return tabs
}
