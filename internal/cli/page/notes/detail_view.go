package notes

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	apimodel "github.com/tugkanmeral/the-host-go/internal/models/api"
)

var (
	detailTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5C2E7"))
	detailMetaStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Faint(true)
	detailBodyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CDD6F4"))

	deleteConfirmBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#F38BA8")).
					Padding(1, 2)
	deleteConfirmTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5C2E7"))
	deleteConfirmHintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CDD6F4"))
)

// FormatDeleteConfirmDialog renders the foreground for the delete-confirmation overlay.
func FormatDeleteConfirmDialog() string {
	title := deleteConfirmTitleStyle.Render("Delete this note?")
	opts := deleteConfirmHintStyle.Render("[y] Yes    [n] No")
	inner := lipgloss.JoinVertical(lipgloss.Left, title, "", opts)
	return deleteConfirmBorderStyle.Render(inner)
}

func FormatNoteDetail(note *apimodel.NoteModel, termWidth int) string {
	if note == nil {
		return EmptyStyle.Render("(no note)")
	}
	innerW := listInnerWidth(termWidth)
	if innerW < 20 {
		innerW = 20
	}

	var blocks []string
	blocks = append(blocks, detailTitleStyle.Width(innerW).Render(note.Title))
	blocks = append(blocks, detailMetaStyle.Render("id  "+note.Id))
	blocks = append(blocks, detailMetaStyle.Render(fmt.Sprintf("created  %s", formatNoteTime(note.CreationDate))))
	if !note.LastUpdateDate.IsZero() {
		blocks = append(blocks, detailMetaStyle.Render(fmt.Sprintf("updated  %s", formatNoteTime(note.LastUpdateDate))))
	}
	blocks = append(blocks, "")

	tagsLine := renderHighlightedTags(note.Tags)
	if strings.TrimSpace(tagsLine) != "" {
		blocks = append(blocks, tagsLine)
		blocks = append(blocks, "")
	}

	body := strings.TrimRight(note.Text, "\n")
	if body != "" {
		blocks = append(blocks, detailBodyStyle.Width(innerW).Render(body))
	}

	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}

func formatNoteTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Local().Format("2006-01-02 15:04")
}
