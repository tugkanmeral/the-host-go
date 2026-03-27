package reminders

import "github.com/charmbracelet/lipgloss"

// PlaceholderBody is shown until the reminders feature is implemented.
func PlaceholderBody(faint lipgloss.Style) string {
	return faint.Render("Nothing here yet. Reminders will appear in this area.")
}
