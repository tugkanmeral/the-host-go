package page

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tugkanmeral/the-host-go/internal/cli/appsvc"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/notes"
	apimodel "github.com/tugkanmeral/the-host-go/internal/models/api"
)

type model struct {
	svc  *appsvc.AppServices
	step Step

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

	listText  string
	listItems []apimodel.NoteListingItemModel
	listTotal int
	listSkip  int
	listTake  int
	listVP    viewport.Model

	listAwaitTakeDigit bool

	noteDetail *apimodel.NoteModel
	detailVP   viewport.Model

	detailDeleteConfirm bool
	detailDeleteLoading bool

	info    string
	errLine string

	// infoReturnToList: StepInfo after successful note update — Enter loads the notes list.
	infoReturnToList bool

	width  int
	height int

	rootMenuCursor int
	menuCursor     int
}

func newModel(svc *appsvc.AppServices) model {
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
		step:       StepLoginUser,
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
		listVP:     viewport.New(80, notes.ListScrollViewportHeight(24)),
		detailVP:   viewport.New(80, notes.DetailScrollViewportHeight(24)),
		listTake:   notes.DefaultListTake,
		width:      80,
		height:     24,
	}
}

// NewModel builds the Bubble Tea root model for the Host TUI.
func NewModel(svc *appsvc.AppServices) tea.Model {
	return newModel(svc)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return tea.ClearScreen() },
		textinput.Blink,
	)
}
