package page

func (m *model) resetListView() {
	m.listText = ""
	m.listItems = nil
	m.listTotal = 0
	m.listSkip = 0
	m.listLoading = false
	m.listAwaitTakeDigit = false
	m.listRequestID++
	m.listSearchTI.SetValue("")
	m.listSearchTI.Blur()
	m.listSearchActive = false
	m.listSearchApplied = ""
	m.listVP.SetContent("")
}
