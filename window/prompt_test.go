package window

import (
	"testing"
)

func Test_loadPrompts(t *testing.T) {
	info := loadPrompts("./embed/prompts-en.csv")
	if len(info) == 0 {
		t.Error("len(info)==0")
	}
	info2 := loadPrompts("./embed/prompts-zh.json")
	if len(info2) == 0 {
		t.Error("len(info2)==0")
	}
}
