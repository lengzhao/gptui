package window

import (
	"testing"
)

func Test_loadPrompts(t *testing.T) {
	info := loadPrompts("../prompts/prompts.csv")
	if len(info) == 0 {
		t.Error("len(info)==0")
	}
	info2 := loadPrompts("../prompts/prompts-zh.json")
	if len(info2) == 0 {
		t.Error("len(info2)==0")
	}
}
