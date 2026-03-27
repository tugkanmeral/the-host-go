package passwords

import "github.com/charmbracelet/lipgloss"

// PlaceholderBody is shown until the passwords feature is implemented.
func PlaceholderBody(faint lipgloss.Style) string {
	return faint.Render("Nothing here yet. Password entries will appear in this area.")
}
