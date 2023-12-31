package chat

import (
	"fmt"
	"strings"

	"github.com/lengzhao/gptui/event"
)

func StartWithEvent() error {
	client, err := newGPTClient()
	if err != nil {
		return err
	}
	c := NewWithClient(client)
	if c == nil {
		return fmt.Errorf("fail to new chat")
	}
	event.RegistEvent(event.EAll, func(key event.EventID, info string) {
		switch key {
		case event.EUserCommit:
			info = strings.TrimSpace(info)
			if len(info) == 0 {
				break
			}
			event.SendEvent(event.EStartChat, info)
			out, err := c.Send(info)
			if err != nil {
				event.SendEvent(event.EChatError, err.Error())
				return
			}
			event.SendEvent(event.EFinishChat, out)
		case event.ESystemPrompt:
			c.SetSystemPrompt(info)
			c.Reset()
		case event.EReset:
			c.Reset()
		}
	})
	event.SendEvent(event.EChatEnable, "")
	return nil
}
