package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tugkanmeral/the-host-go/internal/cli/appsvc"
	"github.com/tugkanmeral/the-host-go/internal/cli/page"
	"github.com/tugkanmeral/the-host-go/internal/service"
)

// Run starts the Bubble Tea TUI using the given services (MongoDB must already be connected).
func Run(authSvc *service.AuthService, noteSvc *service.NoteService) error {
	svc := appsvc.NewAppServices(authSvc, noteSvc)
	p := tea.NewProgram(page.NewModel(svc), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
