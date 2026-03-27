package page

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func navHint(text string) string {
	return navHintStyle.Render(text)
}

// pinFooterBelowContent lays out content + flex + footer like StepListView:
// main block, then "\n\n" (blank line before nav), then nav hint — same as list VP + "\n\n" + navHint.
// flex is nudged until lipgloss.Height(full)==termHeight so bubbles widgets don’t leave a short frame
// (which would make padViewToTerminalHeight append newlines *below* the hint and look “floating”).
func pinFooterBelowContent(content string, footer string, termHeight int) string {
	footerBlock := "\n\n" + footer
	if termHeight <= 0 {
		termHeight = 24
	}
	if termHeight < 6 {
		return content + footerBlock
	}
	hc := lipgloss.Height(content)
	hf := lipgloss.Height(footerBlock)
	flex := termHeight - hc - hf
	if flex < 0 {
		flex = 0
	}
	build := func(f int) string {
		return content + strings.Repeat("\n", f) + footerBlock
	}
	out := build(flex)
	// Grow flex until the block reaches term height (no flex < termHeight cap: bubble widgets
	// may need more spacer rows than the naive termHeight-hc-hf estimate).
	maxFlex := flex + termHeight + 200
	for lipgloss.Height(out) < termHeight && flex < maxFlex {
		flex++
		out = build(flex)
	}
	for flex > 0 && lipgloss.Height(out) > termHeight {
		flex--
		out = build(flex)
	}
	return out
}

// padViewToTerminalHeight pads with blank lines so the frame covers the full terminal height,
// avoiding leftover characters from the scrollback buffer on rows the renderer does not overwrite.
func padViewToTerminalHeight(s string, termHeight int) string {
	if termHeight <= 0 {
		termHeight = 24
	}
	n := lipgloss.Height(s)
	if n >= termHeight {
		return s
	}
	return s + strings.Repeat("\n", termHeight-n)
}
