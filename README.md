# gptui

gpt ui, support openai/azure, enable set http proxy

1. git clone https://github.com/lengzhao/gptui.git
2. go generate ./...
   1. optional: add version info
3. go run main.go
4. go to setting page, set api key, commit
5. restart app

## package

1. go install fyne.io/fyne/v2/cmd/fyne@latest
2. fyne package -os darwin -icon logo.png
