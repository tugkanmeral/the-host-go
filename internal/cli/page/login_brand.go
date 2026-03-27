package page

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const loginBannerASCII = `████████╗██╗  ██╗███████╗    ██╗  ██╗ ██████╗ ███████╗████████╗
╚══██╔══╝██║  ██║██╔════╝    ██║  ██║██╔═══██╗██╔════╝╚══██╔══╝
   ██║   ███████║█████╗      ███████║██║   ██║███████╗   ██║   
   ██║   ██╔══██║██╔══╝      ██╔══██║██║   ██║╚════██║   ██║   
   ██║   ██║  ██║███████╗    ██║  ██║╚██████╔╝███████║   ██║   
   ╚═╝   ╚═╝  ╚═╝╚══════╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝   ╚═╝   `

func loginBranding() string {
	lines := strings.Split(strings.TrimSpace(loginBannerASCII), "\n")
	maxW := 0
	for _, line := range lines {
		if w := lipgloss.Width(line); w > maxW {
			maxW = w
		}
	}
	colored := make([]string, len(lines))
	for i, line := range lines {
		c := loginBannerLineColors[i%len(loginBannerLineColors)]
		colored[i] = lipgloss.NewStyle().Bold(true).Foreground(c).Render(line)
	}
	banner := lipgloss.JoinVertical(lipgloss.Left, colored...)
	spaceLine := "\n"
	tag1 := lipgloss.PlaceHorizontal(maxW, lipgloss.Right, loginTaglinePoweredStyle.Render("powered by Go Lang"))
	tag2 := lipgloss.PlaceHorizontal(maxW, lipgloss.Right, loginTaglineHallucinatedStyle.Render("hallucinated by AI"))
	tag3 := lipgloss.PlaceHorizontal(maxW, lipgloss.Right, loginTaglineDebuggedStyle.Render("debugged by Tuğkan"))
	block := lipgloss.JoinVertical(lipgloss.Left,
		banner,
		spaceLine,
		tag1,
		tag2,
		tag3,
	)
	return lipgloss.NewStyle().MarginTop(1).MarginLeft(2).Render(block)
}
