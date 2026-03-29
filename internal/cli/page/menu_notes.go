package page

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/notes"
)

var notesMenuRows = []struct {
	key  string
	text string
}{
	{"1", "List notes"},
	{"2", "Add note"},
}

func notesMenuRender(m model) string {
	lines := make([]string, len(notesMenuRows))
	for i, row := range notesMenuRows {
		line := fmt.Sprintf("[%s] %s", row.key, row.text)
		if i == m.menuCursor {
			lines[i] = menuItemSelStyle.Render("▸ " + line)
		} else {
			lines[i] = menuItemStyle.Render("  " + line)
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) activateNotesMenuSelection() (tea.Model, tea.Cmd) {
	n := len(notesMenuRows)
	if m.menuCursor < 0 || m.menuCursor >= n {
		return m, nil
	}
	switch m.menuCursor {
	case 0:
		m.listSkip = 0
		m.listTake = notes.NormalizeListTake(m.listTake)
		m.step = StepListLoading
		return m, listCmd(m.svc, m.listSkip, m.listTake)
	case 1:
		m.titleTI.SetValue("")
		m.bodyTA.SetValue("")
		m.tagsTI.SetValue("")
		m.titleTI.Focus()
		m.step = StepAddTitle
		return m, textinput.Blink
	}
	return m, nil
}

func (m model) updateNotesMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		n := len(notesMenuRows)
		switch s {
		case "up", "k":
			m.menuCursor = (m.menuCursor + n - 1) % n
			return m, nil
		case "down", "j":
			m.menuCursor = (m.menuCursor + 1) % n
			return m, nil
		case "enter":
			return m.activateNotesMenuSelection()
		case "1":
			m.menuCursor = 0
			return m.activateNotesMenuSelection()
		case "2":
			m.menuCursor = 1
			return m.activateNotesMenuSelection()
		}
	}
	return m, nil
}
