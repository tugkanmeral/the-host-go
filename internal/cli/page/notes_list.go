package page

func (m *model) resetListView() {
	m.listText = ""
	m.listItems = nil
	m.listTotal = 0
	m.listSkip = 0
	m.listAwaitTakeDigit = false
	m.listVP.SetContent("")
}
