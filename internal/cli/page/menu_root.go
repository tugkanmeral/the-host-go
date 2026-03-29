package page

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var rootMenuRows = []struct {
	key  string
	text string
}{
	{"1", "Note"},
	{"2", "Password"},
	{"3", "Reminder"},
}

func rootMenuRender(m model) string {
	lines := make([]string, len(rootMenuRows))
	for i, row := range rootMenuRows {
		line := fmt.Sprintf("[%s] %s", row.key, row.text)
		if i == m.rootMenuCursor {
			lines[i] = menuItemSelStyle.Render("▸ " + line)
		} else {
			lines[i] = menuItemStyle.Render("  " + line)
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) activateRootMenuSelection() (tea.Model, tea.Cmd) {
	n := len(rootMenuRows)
	if m.rootMenuCursor < 0 || m.rootMenuCursor >= n {
		return m, nil
	}
	switch m.rootMenuCursor {
	case 0:
		m.menuCursor = 0
		m.step = StepNotesMenu
		return m, nil
	case 1:
		m.step = StepPasswordsMenu
		return m, nil
	case 2:
		m.step = StepRemindersMenu
		return m, nil
	}
	return m, nil
}

func (m model) updateRootMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		n := len(rootMenuRows)
		switch s {
		case "up", "k":
			m.rootMenuCursor = (m.rootMenuCursor + n - 1) % n
			return m, nil
		case "down", "j":
			m.rootMenuCursor = (m.rootMenuCursor + 1) % n
			return m, nil
		case "enter":
			return m.activateRootMenuSelection()
		case "1":
			m.rootMenuCursor = 0
			return m.activateRootMenuSelection()
		case "2":
			m.rootMenuCursor = 1
			return m.activateRootMenuSelection()
		case "3":
			m.rootMenuCursor = 2
			return m.activateRootMenuSelection()
		}
	}
	return m, nil
}
