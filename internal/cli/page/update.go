package page

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/notes"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/passwords"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/reminders"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := msg.Width - 4
		if w < 20 {
			w = 20
		}
		m.userTI.Width = w
		m.passTI.Width = w
		m.titleTI.Width = w
		m.tagsTI.Width = w
		m.updIDTI.Width = w
		m.updTitleTI.Width = w
		m.updTagsTI.Width = w
		m.delIDTI.Width = w
		m.bodyTA.SetWidth(w)
		bodyH := msg.Height / 4
		if bodyH < 4 {
			bodyH = 4
		}
		if bodyH > 12 {
			bodyH = 12
		}
		m.bodyTA.SetHeight(bodyH)
		m.updBodyTA.SetWidth(w)
		uh := msg.Height / 5
		if uh < 3 {
			uh = 3
		}
		if uh > 8 {
			uh = 8
		}
		m.updBodyTA.SetHeight(uh)
		if m.step == StepListView {
			m.listText = notes.FormatNotesList(m.listItems, m.listSkip, m.width)
			m.listVP.SetContent(m.listText)
			m.listVP.Width = m.width
			m.listVP.Height = notes.ListScrollViewportHeight(m.height)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+c", "esc":
			switch m.step {
			case StepLoginUser, StepLoginPass, StepRootMenu:
				return m, tea.Quit
			case StepNotesMenu, StepPasswordsMenu, StepRemindersMenu:
				m.menuCursor = 0
				m.step = StepRootMenu
				m.errLine = ""
				return m, nil
			case StepListView:
				m.resetListView()
				m.menuCursor = 0
				m.step = StepNotesMenu
				m.errLine = ""
				return m, nil
			case StepListLoading:
				m.step = StepNotesMenu
				m.errLine = ""
				return m, nil
			case StepAddTitle, StepAddText, StepAddTags,
				StepUpdateID, StepUpdateTitle, StepUpdateText, StepUpdateTags, StepDeleteID:
				m.menuCursor = 0
				m.step = StepNotesMenu
				m.errLine = ""
				return m, nil
			case StepInfo:
				m.menuCursor = 0
				m.step = StepNotesMenu
				m.info = ""
				m.errLine = ""
				return m, nil
			}
		}

	case loginDoneMsg:
		if msg.err != nil {
			m.errLine = msg.err.Error()
			m.step = StepLoginUser
			m.passTI.SetValue("")
			m.userTI.Focus()
			return m, textinput.Blink
		}
		m.errLine = ""
		m.rootMenuCursor = 0
		m.menuCursor = 0
		m.step = StepRootMenu
		return m, func() tea.Msg { return tea.ClearScreen() }

	case listDoneMsg:
		if msg.err != nil {
			m.info = ""
			m.errLine = msg.err.Error()
			m.step = StepInfo
			return m, nil
		}
		m.listItems = msg.items
		m.listTotal = msg.total
		m.listSkip = msg.skip
		m.listTake = msg.take
		m.listText = notes.FormatNotesList(msg.items, msg.skip, m.width)
		m.listVP.SetContent(m.listText)
		m.listVP.Width = m.width
		m.listVP.Height = notes.ListScrollViewportHeight(m.height)
		m.listVP.GotoTop()
		m.step = StepListView
		return m, nil

	case simpleErrMsg:
		m.errLine = msg.err.Error()
		m.step = StepInfo
		return m, nil

	case simpleOkMsg:
		m.errLine = ""
		m.info = "Operation completed successfully."
		m.step = StepInfo
		return m, nil
	}

	switch m.step {
	case StepLoginUser:
		return m.updateLoginUser(msg)
	case StepLoginPass:
		return m.updateLoginPass(msg)
	case StepRootMenu:
		return m.updateRootMenu(msg)
	case StepNotesMenu:
		return m.updateNotesMenu(msg)
	case StepPasswordsMenu:
		return m, nil
	case StepRemindersMenu:
		return m, nil
	case StepListLoading:
		return m, nil
	case StepListView:
		return m.updateListView(msg)
	case StepAddTitle:
		return m.updateAddTitle(msg)
	case StepAddText:
		return m.updateAddText(msg)
	case StepAddTags:
		return m.updateAddTags(msg)
	case StepUpdateID:
		return m.updateUpdID(msg)
	case StepUpdateTitle:
		return m.updateUpdTitle(msg)
	case StepUpdateText:
		return m.updateUpdText(msg)
	case StepUpdateTags:
		return m.updateUpdTags(msg)
	case StepDeleteID:
		return m.updateDeleteID(msg)
	case StepInfo:
		return m.updateInfo(msg)
	}

	return m, nil
}

func (m model) updateLoginUser(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			u := strings.TrimSpace(m.userTI.Value())
			if u == "" {
				m.errLine = "Username is required."
				return m, nil
			}
			m.errLine = ""
			m.userTI.Blur()
			m.passTI.Focus()
			m.step = StepLoginPass
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.userTI, cmd = m.userTI.Update(msg)
	return m, cmd
}

func (m model) updateLoginPass(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			user := strings.TrimSpace(m.userTI.Value())
			pass := m.passTI.Value()
			if pass == "" {
				m.errLine = "Password is required."
				return m, nil
			}
			m.errLine = ""
			m.passTI.Blur()
			return m, loginCmd(m.svc, user, pass)
		}
	}
	var cmd tea.Cmd
	m.passTI, cmd = m.passTI.Update(msg)
	return m, cmd
}

func (m model) reloadListWithTake(newTake int) (tea.Model, tea.Cmd) {
	newTake = notes.NormalizeListTake(newTake)
	skip := m.listSkip
	if newTake != m.listTake {
		skip = 0
	}
	m.listSkip = skip
	m.listTake = newTake
	m.step = StepListLoading
	return m, listCmd(m.svc, skip, newTake)
}

func (m model) listLoadPrevPage() (tea.Model, tea.Cmd) {
	take := notes.NormalizeListTake(m.listTake)
	if m.listSkip <= 0 {
		return m, nil
	}
	ns := m.listSkip - take
	if ns < 0 {
		ns = 0
	}
	m.listSkip = ns
	m.step = StepListLoading
	return m, listCmd(m.svc, ns, take)
}

func (m model) listLoadNextPage() (tea.Model, tea.Cmd) {
	take := notes.NormalizeListTake(m.listTake)
	if m.listSkip+take >= m.listTotal {
		return m, nil
	}
	ns := m.listSkip + take
	m.listSkip = ns
	m.step = StepListLoading
	return m, listCmd(m.svc, ns, take)
}

func listTakeFromDigitKey(msg tea.KeyMsg) (take int, ok bool) {
	if msg.Type != tea.KeyRunes || msg.Paste || len(msg.Runes) != 1 {
		return 0, false
	}
	r := msg.Runes[0]
	if r < '0' || r > '9' {
		return 0, false
	}
	if r == '0' {
		return 10, true
	}
	return int(r - '0'), true
}

func (m model) updateListView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "enter":
			m.resetListView()
			m.menuCursor = 0
			m.step = StepNotesMenu
			return m, nil
		}
		if take, ok := listTakeFromDigitKey(msg); ok {
			return m.reloadListWithTake(take)
		}
		switch s {
		case "[", "left", "p":
			return m.listLoadPrevPage()
		case "]", "right", "n":
			return m.listLoadNextPage()
		}
		var cmd tea.Cmd
		m.listVP, cmd = m.listVP.Update(msg)
		return m, cmd
	case tea.MouseMsg:
		var cmd tea.Cmd
		m.listVP, cmd = m.listVP.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) updateAddTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if strings.TrimSpace(m.titleTI.Value()) == "" {
				m.errLine = "Title is required."
				return m, nil
			}
			m.errLine = ""
			m.titleTI.Blur()
			m.bodyTA.Focus()
			m.step = StepAddText
			return m, textarea.Blink
		}
	}
	var cmd tea.Cmd
	m.titleTI, cmd = m.titleTI.Update(msg)
	return m, cmd
}

func (m model) updateAddText(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+e" {
			if strings.TrimSpace(m.bodyTA.Value()) == "" {
				m.errLine = "Text is required."
				return m, nil
			}
			m.errLine = ""
			m.bodyTA.Blur()
			m.tagsTI.Focus()
			m.step = StepAddTags
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.bodyTA, cmd = m.bodyTA.Update(msg)
	return m, cmd
}

func (m model) updateAddTags(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			tags := parseTags(m.tagsTI.Value())
			m.tagsTI.Blur()
			title := strings.TrimSpace(m.titleTI.Value())
			text := m.bodyTA.Value()
			return m, addCmd(m.svc, title, text, tags)
		}
	}
	var cmd tea.Cmd
	m.tagsTI, cmd = m.tagsTI.Update(msg)
	return m, cmd
}

func (m model) updateUpdID(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if strings.TrimSpace(m.updIDTI.Value()) == "" {
				m.errLine = "Note ID is required."
				return m, nil
			}
			m.errLine = ""
			m.updIDTI.Blur()
			m.updTitleTI.Focus()
			m.step = StepUpdateTitle
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.updIDTI, cmd = m.updIDTI.Update(msg)
	return m, cmd
}

func (m model) updateUpdTitle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.updTitleTI.Blur()
			m.updBodyTA.Focus()
			m.step = StepUpdateText
			return m, textarea.Blink
		}
	}
	var cmd tea.Cmd
	m.updTitleTI, cmd = m.updTitleTI.Update(msg)
	return m, cmd
}

func (m model) updateUpdText(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+e" {
			m.updBodyTA.Blur()
			m.updTagsTI.Focus()
			m.step = StepUpdateTags
			return m, textinput.Blink
		}
	}
	var cmd tea.Cmd
	m.updBodyTA, cmd = m.updBodyTA.Update(msg)
	return m, cmd
}

func (m model) updateUpdTags(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			id := strings.TrimSpace(m.updIDTI.Value())
			title := strings.TrimSpace(m.updTitleTI.Value())
			text := strings.TrimSpace(m.updBodyTA.Value())
			tags := parseTags(m.updTagsTI.Value())
			if title == "" && text == "" && len(tags) == 0 {
				m.errLine = "Provide at least one of: title, text, or tags."
				return m, nil
			}
			m.updTagsTI.Blur()
			return m, updateCmd(m.svc, id, title, text, tags)
		}
	}
	var cmd tea.Cmd
	m.updTagsTI, cmd = m.updTagsTI.Update(msg)
	return m, cmd
}

func (m model) updateDeleteID(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			id := strings.TrimSpace(m.delIDTI.Value())
			if id == "" {
				m.errLine = "Note ID is required."
				return m, nil
			}
			m.errLine = ""
			m.delIDTI.Blur()
			return m, deleteCmd(m.svc, id)
		}
	}
	var cmd tea.Cmd
	m.delIDTI, cmd = m.delIDTI.Update(msg)
	return m, cmd
}

func (m model) updateInfo(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.menuCursor = 0
		m.step = StepNotesMenu
		m.info = ""
		m.errLine = ""
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	var v string
	switch m.step {
	case StepLoginUser:
		s := loginBranding() + "\n\n" + signInTitleStyle.Render("Sign in") + "\n\nUsername:\n" + m.userTI.View()
		if m.errLine != "" {
			s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		}
		v = pinFooterBelowContent(s, navHint("Enter: continue · Esc or Ctrl+Q: quit"), m.height)
	case StepLoginPass:
		s := loginBranding() + "\n\n" + signInTitleStyle.Render("Sign in") + "\n\nPassword:\n" + m.passTI.View()
		if m.errLine != "" {
			s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		}
		v = pinFooterBelowContent(s, navHint("Enter: login · Esc or Ctrl+Q: quit"), m.height)
	case StepRootMenu:
		s := loginBranding() + "\n\n" + titleStyle.Render("Main menu") + "\n\n" + rootMenuRender(m)
		v = pinFooterBelowContent(s, navHint("↑↓ j k · Enter · 1–3 · Esc / Ctrl+Q: quit app"), m.height)
	case StepNotesMenu:
		s := loginBranding() + "\n\n" + titleStyle.Render("Note menu") + "\n\n" + notesMenuRender(m)
		v = pinFooterBelowContent(s, navHint("↑↓ j k · Enter · 1–2 · Esc: back · Ctrl+Q: quit"), m.height)
	case StepPasswordsMenu:
		s := loginBranding() + "\n\n" + titleStyle.Render("Password menu") + "\n\n" + passwords.PlaceholderBody(sectionPlaceholderStyle)
		v = pinFooterBelowContent(s, navHint("Esc: back to main menu · Ctrl+Q: quit"), m.height)
	case StepRemindersMenu:
		s := loginBranding() + "\n\n" + titleStyle.Render("Reminder menu") + "\n\n" + reminders.PlaceholderBody(sectionPlaceholderStyle)
		v = pinFooterBelowContent(s, navHint("Esc: back to main menu · Ctrl+Q: quit"), m.height)
	case StepListLoading:
		v = notes.ScreenTitleStyle.Render("Notes") + "\n\n" + notes.LoadingStyle.Render("Loading notes…")
	case StepListView:
		banner := notes.ListViewPagingBanner(m.listItems, m.listTotal, m.listSkip, m.listTake)
		v = notes.ScreenTitleStyle.Render("Notes") + "\n" + banner + "\n\n" + m.listVP.View() + "\n\n" +
			navHint("0–9 = page size (0→10) · ←p →n [] · ↑↓jk scroll · Enter/Esc: Note menu · Ctrl+Q: quit")
	case StepAddTitle:
		s := titleStyle.Render("New note — title") + "\n\n" + m.titleTI.View()
		v = pinFooterBelowContent(s, navHint("Enter: continue · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepAddText:
		s := titleStyle.Render("New note — body") + "\n\n" + m.bodyTA.View()
		v = pinFooterBelowContent(s, navHint("Ctrl+E: continue · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepAddTags:
		s := titleStyle.Render("New note — tags (optional)") + "\n\n" + m.tagsTI.View()
		v = pinFooterBelowContent(s, navHint("Enter: save · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepUpdateID:
		s := titleStyle.Render("Update note — id") + "\n\n" + m.updIDTI.View()
		v = pinFooterBelowContent(s, navHint("Enter: continue · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepUpdateTitle:
		s := titleStyle.Render("Update — title") + "\n\n" + m.updTitleTI.View()
		v = pinFooterBelowContent(s, navHint("Enter: continue · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepUpdateText:
		s := titleStyle.Render("Update — text") + "\n\n" + m.updBodyTA.View()
		v = pinFooterBelowContent(s, navHint("Ctrl+E: continue · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepUpdateTags:
		s := titleStyle.Render("Update — tags") + "\n\n" + m.updTagsTI.View()
		if m.errLine != "" {
			s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		}
		v = pinFooterBelowContent(s, navHint("Enter: submit · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepDeleteID:
		s := titleStyle.Render("Delete note") + "\n\n" + m.delIDTI.View()
		v = pinFooterBelowContent(s, navHint("Enter: delete · Esc: menu · Ctrl+Q: quit"), m.height)
	case StepInfo:
		var s string
		if m.errLine != "" {
			s = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		} else {
			s = m.info
		}
		v = pinFooterBelowContent(s, navHint("Any key: Note menu · Ctrl+Q: quit"), m.height)
	default:
		v = ""
	}
	// List screens use manual viewport height; pad clears scrollback below the nav line.
	// pinFooterBelowContent already sizes to term height — extra pad would add blank lines *under* the hint.
	switch m.step {
	case StepListView, StepListLoading:
		return padViewToTerminalHeight(v, m.height)
	default:
		if v == "" {
			return padViewToTerminalHeight(v, m.height)
		}
		return v
	}
}
