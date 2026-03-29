package notes

import "github.com/charmbracelet/lipgloss"

var (
	cardDividerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true)

	ScreenTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CBA6F7"))
	PagingStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)
	EmptyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true).Italic(true)
	LoadingStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)

	noteBorderColors = [...]lipgloss.Color{
		lipgloss.Color("#CBA6F7"),
		lipgloss.Color("#89B4FA"),
		lipgloss.Color("#A6E3A1"),
		lipgloss.Color("#F9E2AF"),
	}
	noteTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5C2E7"))
	noteIDStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)
	noteBodyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CDD6F4"))

	noteTagsLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)
	tagChipFG          = lipgloss.Color("#11111B")
	tagChipBGColors    = [...]lipgloss.Color{
		lipgloss.Color("#94E2D5"),
		lipgloss.Color("#89DCEB"),
		lipgloss.Color("#CBA6F7"),
		lipgloss.Color("#F9E2AF"),
	}

	noteHotkeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#11111B")).
			Background(lipgloss.Color("#F9E2AF")).
			Padding(0, 1)
)
