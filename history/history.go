package history

type HistoryItem struct {
	Role string
	Text string
}
type History struct {
	items []HistoryItem
}

func New() *History {
	return &History{}
}

func (m *History) Add(role, text string) {
	m.items = append(m.items, HistoryItem{
		Role: role,
		Text: text,
	})
}

func (m *History) Get(limit int) []HistoryItem {
	if len(m.items) == 0 {
		return nil
	}
	if limit > len(m.items) {
		return m.items
	}
	return m.items[len(m.items)-limit : len(m.items)]
}
