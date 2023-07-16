package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/lengzhao/conf"
	_ "github.com/lengzhao/conf/autoload"
	_ "github.com/lengzhao/font/autoload"
	"github.com/lengzhao/gpt/chat"
	"github.com/lengzhao/gpt/window"
)

func main() {
	a := app.New()
	w := a.NewWindow("gpt")
	w.Resize(fyne.NewSize(float32(conf.GetFloat("windWidth", 1000)),
		float32(conf.GetFloat("windHeight", 700))))

	w.SetContent(makeTabs())
	switch conf.Get("theme", "default") {
	case "dark":
		a.Settings().SetTheme(theme.DarkTheme())
	case "light":
		a.Settings().SetTheme(theme.LightTheme())
	default:
		a.Settings().SetTheme(theme.DefaultTheme())
	}
	chat.StartWithEvent()

	w.ShowAndRun()
}

func makeTabs() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Chat", theme.AccountIcon(), window.MakeChatWindow()),
		container.NewTabItemWithIcon("Role", theme.DocumentPrintIcon(), window.MakeRoleWindow()),
		container.NewTabItemWithIcon("Setting", theme.SettingsIcon(), window.MakeSettingWindow()),
	)
	return tabs
}
