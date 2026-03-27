package page

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle            = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	signInTitleStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A78BFA"))
	// loginTaglineStyle: same “washed / dim” treatment as section placeholders (Faint), with hue from each derived style.
	loginTaglineStyle = lipgloss.NewStyle().Faint(true)

	loginTaglinePoweredStyle = lipgloss.NewStyle().
				Inherit(loginTaglineStyle).
				Foreground(lipgloss.Color("#90CAF9"))
	loginTaglineHallucinatedStyle = lipgloss.NewStyle().
					Inherit(loginTaglineStyle).
					Foreground(lipgloss.Color("#B39DDB"))
	loginTaglineDebuggedStyle = lipgloss.NewStyle().
					Inherit(loginTaglineStyle).
					Foreground(lipgloss.Color("#A5D6A7"))
	loginBannerLineColors = [...]lipgloss.Color{
		lipgloss.Color("#CBA6F7"), // mauve
		lipgloss.Color("#F5C2E7"), // pink
		lipgloss.Color("#89B4FA"), // blue
		lipgloss.Color("#B4BEFE"), // lavender
		lipgloss.Color("#A6E3A1"), // green
		lipgloss.Color("#A78BFA"), // violet
	}
	navHintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Faint(true)
	menuItemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	menuItemSelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5C2E7"))

	sectionPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)
)
