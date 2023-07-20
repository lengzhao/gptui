package event

import (
	"sync"

	"github.com/lengzhao/conf"
)

//go:generate stringer -type=EventID
type EventID int
type EventCb func(key EventID, info string)

const (
	EHistory EventID = iota
	EStartChat
	EFinishChat
	EUserCommit
	EChatError
	EReset
	ERetry
	ESystemPrompt
	ELoadPrompt
	EChatEnable
	EError
	EAll
)

var eventReg map[EventID][]EventCb

func init() {
	eventReg = make(map[EventID][]EventCb)
}

/*
if key == EAll, receive all event
*/
func RegistEvent(key EventID, cb EventCb) {
	list := eventReg[key]
	list = append(list, cb)
	eventReg[key] = list
}

type eventItem struct {
	key  EventID
	info string
}

var events []eventItem
var mux sync.Mutex

func SendEvent(key EventID, info string) {
	if key == EAll {
		return
	}
	mux.Lock()
	events = append(events, eventItem{
		key:  key,
		info: info,
	})
	if len(events) > 1 {
		mux.Unlock()
		return
	}
	mux.Unlock()
	for i := 0; i < conf.GetInt("MaxEventPerBatch", 500); i++ {
		e := popEvent()
		if e == nil {
			return
		}
		list2 := eventReg[EAll]
		for _, it := range list2 {
			it(e.key, e.info)
		}
		list := eventReg[e.key]
		for _, it := range list {
			it(e.key, e.info)
		}
	}
}

func popEvent() *eventItem {
	if len(events) == 0 {
		return nil
	}
	mux.Lock()
	e := events[0]
	events = events[1:]
	mux.Unlock()
	return &e
}
