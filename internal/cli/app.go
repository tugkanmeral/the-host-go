package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tugkanmeral/the-host-go/internal/service"
)

type step int

const (
	stepLoginUser step = iota
	stepLoginPass
	stepMenu
	stepListLoading
	stepListView
	stepAddTitle
	stepAddText
	stepAddTags
	stepUpdateID
	stepUpdateTitle
	stepUpdateText
	stepUpdateTags
	stepDeleteID
	stepInfo
)

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))

type model struct {
	svc  *AppServices
	step step

	userTI textinput.Model
	passTI textinput.Model

	titleTI textinput.Model
	bodyTA  textarea.Model
	tagsTI  textinput.Model

	updIDTI    textinput.Model
	updTitleTI textinput.Model
	updBodyTA  textarea.Model
	updTagsTI  textinput.Model

	delIDTI textinput.Model

	listText string
	info     string
	errLine  string

	width  int
	height int
}

func newModel(svc *AppServices) model {
	ui := textinput.New()
	ui.Placeholder = "username"
	ui.Focus()
	ui.CharLimit = 64
	ui.Width = 40

	pi := textinput.New()
	pi.Placeholder = "password"
	pi.EchoMode = textinput.EchoPassword
	pi.CharLimit = 128
	pi.Width = 40

	tti := textinput.New()
	tti.Placeholder = "title"
	tti.CharLimit = 200
	tti.Width = 50

	ta := textarea.New()
	ta.Placeholder = "note body"
	ta.SetWidth(50)
	ta.SetHeight(6)
	ta.CharLimit = 8000

	tagi := textinput.New()
	tagi.Placeholder = "tags (comma-separated, optional)"
	tagi.CharLimit = 500
	tagi.Width = 50

	uid := textinput.New()
	uid.Placeholder = "note id"
	uid.CharLimit = 32
	uid.Width = 40

	ut := textinput.New()
	ut.Placeholder = "new title (empty = skip)"
	ut.CharLimit = 200
	ut.Width = 50

	uta := textarea.New()
	uta.Placeholder = "new text (empty = skip)"
	uta.SetWidth(50)
	uta.SetHeight(5)
	uta.CharLimit = 8000

	utg := textinput.New()
	utg.Placeholder = "new tags (comma-separated, empty = skip)"
	utg.CharLimit = 500
	utg.Width = 50

	did := textinput.New()
	did.Placeholder = "note id to delete"
	did.CharLimit = 32
	did.Width = 40

	return model{
		svc:        svc,
		step:       stepLoginUser,
		userTI:     ui,
		passTI:     pi,
		titleTI:    tti,
		bodyTA:     ta,
		tagsTI:     tagi,
		updIDTI:    uid,
		updTitleTI: ut,
		updBodyTA:  uta,
		updTagsTI:  utg,
		delIDTI:    did,
		width:      80,
		height:     24,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type loginDoneMsg struct{ err error }
type listDoneMsg struct {
	text string
	err  error
}
type simpleErrMsg struct{ err error }
type simpleOkMsg struct{}

func loginCmd(svc *AppServices, user, pass string) tea.Cmd {
	return func() tea.Msg {
		return loginDoneMsg{err: svc.Login(context.Background(), user, pass)}
	}
}

func listCmd(svc *AppServices) tea.Cmd {
	return func() tea.Msg {
		items, total, err := svc.ListNotes(context.Background(), 0, 50)
		if err != nil {
			return listDoneMsg{err: err}
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("Total: %d\n\n", total))
		for i, it := range items {
			b.WriteString(fmt.Sprintf("%d) [%s] %s\n", i+1, it.Id, it.Title))
			if it.Text != "" {
				preview := it.Text
				if len(preview) > 120 {
					preview = preview[:120] + "…"
				}
				b.WriteString("   ")
				b.WriteString(preview)
				b.WriteString("\n")
			}
			if len(it.Tags) > 0 {
				b.WriteString(fmt.Sprintf("   tags: %s\n", strings.Join(it.Tags, ", ")))
			}
			b.WriteString("\n")
		}
		if len(items) == 0 {
			b.WriteString("(no notes)\n")
		}
		return listDoneMsg{text: b.String()}
	}
}

func addCmd(svc *AppServices, title, text string, tags []string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.AddNote(context.Background(), title, text, tags); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}

func updateCmd(svc *AppServices, id, title, text string, tags []string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.UpdateNote(context.Background(), id, title, text, tags); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}

func deleteCmd(svc *AppServices, id string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.DeleteNote(context.Background(), id); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}

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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.step == stepMenu {
				return m, tea.Quit
			}
			if m.step == stepLoginUser || m.step == stepLoginPass {
				return m, tea.Quit
			}
			m.step = stepMenu
			m.errLine = ""
			return m, nil
		}

	case loginDoneMsg:
		if msg.err != nil {
			m.errLine = msg.err.Error()
			m.step = stepLoginUser
			m.passTI.SetValue("")
			m.userTI.Focus()
			return m, textinput.Blink
		}
		m.errLine = ""
		m.step = stepMenu
		return m, nil

	case listDoneMsg:
		if msg.err != nil {
			m.info = ""
			m.errLine = msg.err.Error()
			m.step = stepInfo
			return m, nil
		}
		m.listText = msg.text
		m.step = stepListView
		return m, nil

	case simpleErrMsg:
		m.errLine = msg.err.Error()
		m.step = stepInfo
		return m, nil

	case simpleOkMsg:
		m.errLine = ""
		m.info = "Operation completed successfully."
		m.step = stepInfo
		return m, nil
	}

	switch m.step {
	case stepLoginUser:
		return m.updateLoginUser(msg)
	case stepLoginPass:
		return m.updateLoginPass(msg)
	case stepMenu:
		return m.updateMenu(msg)
	case stepListLoading:
		return m, nil
	case stepListView:
		return m.updateListView(msg)
	case stepAddTitle:
		return m.updateAddTitle(msg)
	case stepAddText:
		return m.updateAddText(msg)
	case stepAddTags:
		return m.updateAddTags(msg)
	case stepUpdateID:
		return m.updateUpdID(msg)
	case stepUpdateTitle:
		return m.updateUpdTitle(msg)
	case stepUpdateText:
		return m.updateUpdText(msg)
	case stepUpdateTags:
		return m.updateUpdTags(msg)
	case stepDeleteID:
		return m.updateDeleteID(msg)
	case stepInfo:
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
			m.step = stepLoginPass
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

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.step = stepListLoading
			return m, listCmd(m.svc)
		case "2":
			m.titleTI.SetValue("")
			m.bodyTA.SetValue("")
			m.tagsTI.SetValue("")
			m.titleTI.Focus()
			m.step = stepAddTitle
			return m, textinput.Blink
		case "3":
			m.updIDTI.SetValue("")
			m.updTitleTI.SetValue("")
			m.updBodyTA.SetValue("")
			m.updTagsTI.SetValue("")
			m.updIDTI.Focus()
			m.step = stepUpdateID
			return m, textinput.Blink
		case "4":
			m.delIDTI.SetValue("")
			m.delIDTI.Focus()
			m.step = stepDeleteID
			return m, textinput.Blink
		case "q", "Q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateListView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyMsg); ok {
		m.step = stepMenu
		m.listText = ""
		return m, nil
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
			m.step = stepAddText
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
			m.step = stepAddTags
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
			m.step = stepUpdateTitle
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
			m.step = stepUpdateText
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
			m.step = stepUpdateTags
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
		m.step = stepMenu
		m.info = ""
		m.errLine = ""
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	switch m.step {
	case stepLoginUser:
		s := titleStyle.Render("The Host — Sign in") + "\n\nUsername:\n" + m.userTI.View()
		if m.errLine != "" {
			s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		}
		s += "\n\nEnter: continue · Esc: quit"
		return s
	case stepLoginPass:
		s := titleStyle.Render("The Host — Sign in") + "\n\nPassword:\n" + m.passTI.View()
		if m.errLine != "" {
			s += "\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine)
		}
		s += "\n\nEnter: login · Esc: quit"
		return s
	case stepMenu:
		menu := `[1] List notes
[2] Add note
[3] Update note
[4] Delete note
[q] Quit`
		return titleStyle.Render("Main menu") + "\n\n" + menu + "\n\nPress a key."
	case stepListLoading:
		return "Loading notes…"
	case stepListView:
		return titleStyle.Render("Your notes") + "\n\n" + m.listText + "\n\nPress any key to return."
	case stepAddTitle:
		h := "Enter: continue · Escape: back to menu"
		return titleStyle.Render("New note — title") + "\n\n" + m.titleTI.View() + "\n\n" + h
	case stepAddText:
		h := "Ctrl+E: continue to tags · Escape: menu"
		return titleStyle.Render("New note — body") + "\n\n" + m.bodyTA.View() + "\n\n" + h
	case stepAddTags:
		h := "Enter: save · Escape: menu"
		return titleStyle.Render("New note — tags (optional)") + "\n\n" + m.tagsTI.View() + "\n\n" + h
	case stepUpdateID:
		return titleStyle.Render("Update note — id") + "\n\n" + m.updIDTI.View() + "\n\nEnter: continue · Esc: menu"
	case stepUpdateTitle:
		return titleStyle.Render("Update — title") + "\n\n" + m.updTitleTI.View() + "\n\nEnter: continue · Esc: menu"
	case stepUpdateText:
		return titleStyle.Render("Update — text") + "\n\n" + m.updBodyTA.View() + "\n\nCtrl+E: continue · Esc: menu"
	case stepUpdateTags:
		if m.errLine != "" {
			return titleStyle.Render("Update — tags") + "\n\n" + m.updTagsTI.View() + "\n\n" +
				lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine) + "\n\nEnter: submit · Esc: menu"
		}
		return titleStyle.Render("Update — tags") + "\n\n" + m.updTagsTI.View() + "\n\nEnter: submit · Esc: menu"
	case stepDeleteID:
		return titleStyle.Render("Delete note") + "\n\n" + m.delIDTI.View() + "\n\nEnter: delete · Esc: menu"
	case stepInfo:
		if m.errLine != "" {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.errLine) + "\n\nPress any key for main menu."
		}
		return m.info + "\n\nPress any key for main menu."
	default:
		return ""
	}
}

func parseTags(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// Run starts the Bubble Tea TUI using the given services (MongoDB must already be connected).
func Run(authSvc *service.AuthService, noteSvc *service.NoteService) error {
	svc := NewAppServices(authSvc, noteSvc)
	p := tea.NewProgram(newModel(svc), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
