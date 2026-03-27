package notes

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	apimodel "github.com/tugkanmeral/the-host-go/internal/models/api"
)

func NormalizeListTake(t int) int {
	if t < 1 {
		return DefaultListTake
	}
	if t > MaxListTake {
		return MaxListTake
	}
	return t
}

func listTotalPages(total, take int) int {
	take = NormalizeListTake(take)
	if total <= 0 {
		return 1
	}
	return (total + take - 1) / take
}

func listCurrentPage(skip, take int) int {
	take = NormalizeListTake(take)
	if take < 1 {
		return 1
	}
	return skip/take + 1
}

func ListViewPagingBanner(items []apimodel.NoteListingItemModel, total, skip, take int) string {
	return PagingStyle.Render(formatListPagingSummary(items, total, skip, take))
}

func formatListPagingSummary(items []apimodel.NoteListingItemModel, total, skip, take int) string {
	take = NormalizeListTake(take)
	tp := listTotalPages(total, take)
	cp := listCurrentPage(skip, take)
	n := len(items)
	from, to := 0, 0
	if n > 0 {
		from = skip + 1
		to = skip + n
	}
	return fmt.Sprintf(
		"Page %d/%d · items %d–%d of %d · page size %d · visible on this page: %d",
		cp, tp, from, to, total, take, n,
	)
}

func ListScrollViewportHeight(termHeight int) int {
	h := termHeight - ListViewFrameLines
	if h < 4 {
		h = 4
	}
	return h
}

func renderHighlightedTags(tags []string) string {
	var chips []string
	n := 0
	for _, raw := range tags {
		t := strings.TrimSpace(raw)
		if t == "" {
			continue
		}
		bg := tagChipBGColors[n%len(tagChipBGColors)]
		n++
		chips = append(chips, lipgloss.NewStyle().
			Foreground(tagChipFG).
			Background(bg).
			Padding(0, 1).
			Render(t))
	}
	if len(chips) == 0 {
		return ""
	}
	row := []string{noteTagsLabelStyle.Render("tags"), "  "}
	for i, c := range chips {
		if i > 0 {
			row = append(row, " ")
		}
		row = append(row, c)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, row...)
}

func listInnerWidth(termWidth int) int {
	w := termWidth - 8
	if w < 24 {
		w = 24
	}
	if w > 82 {
		w = 82
	}
	return w
}

func cardSeparator(innerW int) string {
	dashW := innerW - 6
	if dashW < 6 {
		dashW = 6
	}
	line := strings.Repeat("─", dashW)
	return cardDividerStyle.MarginLeft(2).MarginRight(2).Render(line)
}

func renderNoteCard(displayNo int, it apimodel.NoteListingItemModel, termWidth int, idx int) string {
	innerW := listInnerWidth(termWidth)
	border := noteBorderColors[idx%len(noteBorderColors)]
	titleLine := noteTitleStyle.Width(innerW).Render(fmt.Sprintf("%d - %s", displayNo, it.Title))
	idLine := noteIDStyle.Render(it.Id)

	body := strings.TrimSpace(it.Text)
	tagsLine := renderHighlightedTags(it.Tags)
	hasBody := body != ""
	hasTags := tagsLine != ""

	var lines []string
	lines = append(lines, titleLine, idLine)

	if hasBody || hasTags {
		lines = append(lines, cardSeparator(innerW))
	}
	if hasBody {
		runes := []rune(body)
		maxR := innerW * 4
		if len(runes) > maxR {
			body = string(runes[:maxR]) + "…"
		}
		lines = append(lines, noteBodyStyle.Width(innerW).Render(body))
	}
	if hasBody && hasTags {
		lines = append(lines, cardSeparator(innerW))
	}
	if hasTags {
		lines = append(lines, tagsLine)
	}

	core := lipgloss.JoinVertical(lipgloss.Left, lines...)
	maxOuter := termWidth - 2
	if maxOuter < 28 {
		maxOuter = 28
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(0, 1).
		MaxWidth(maxOuter).
		Render(core)
}

func FormatNotesList(items []apimodel.NoteListingItemModel, skip, termWidth int) string {
	parts := []string{}
	if len(items) == 0 {
		parts = append(parts, EmptyStyle.Render("(no notes on this page)"))
		return lipgloss.JoinVertical(lipgloss.Left, parts...)
	}
	for i, it := range items {
		displayNo := skip + i + 1
		parts = append(parts, renderNoteCard(displayNo, it, termWidth, i))
		if i < len(items)-1 {
			parts = append(parts, "")
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
