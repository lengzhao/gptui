package history

import "testing"

func TestHistory_Add(t *testing.T) {
	h := New()
	h.Add("role1", "text1")
	h.Add("role2", "text2")
	h.Add("role3", "text3")
	h.Add("role4", "text4")
	h.Add("role5", "text5")
	h.Add("role6", "text6")
	h.Add("role7", "text7")
	h.Add("role8", "text8")
	h.Add("role9", "text9")
	h.Add("role10", "text10")
	h.Add("role11", "text11")
}

func TestHistory_Get(t *testing.T) {
	h := New()
	rst := h.Get(5)
	if rst != nil {
		t.Error(rst)
	}
	h.Add("role1", "text1")
	rst1 := h.Get(5)
	if len(rst1) != 1 {
		t.Error(rst1)
	}
	h.Add("role2", "text2")
	h.Add("role3", "text3")
	h.Add("role4", "text4")
	h.Add("role5", "text5")
	h.Add("role6", "text6")
	h.Add("role7", "text7")
	h.Add("role8", "text8")
	h.Add("role9", "text9")
	h.Add("role10", "text10")
	h.Add("role11", "text11")
	rst2 := h.Get(5)
	if len(rst2) != 5 {
		t.Error(rst2)
	}
	if rst2[0].Role != "role7" {
		t.Error(rst2[0].Role)
	}
	if rst2[4].Role != "role11" {
		t.Error(rst2[4].Role)
	}
}
