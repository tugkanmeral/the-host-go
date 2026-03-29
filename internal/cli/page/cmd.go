package page

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tugkanmeral/the-host-go/internal/cli/appsvc"
	"github.com/tugkanmeral/the-host-go/internal/cli/page/notes"
	apimodel "github.com/tugkanmeral/the-host-go/internal/models/api"
)

type loginDoneMsg struct{ err error }
type listDoneMsg struct {
	items []apimodel.NoteListingItemModel
	total int
	skip  int
	take  int
	err   error
}
type simpleErrMsg struct{ err error }
type simpleOkMsg struct{}

type noteDetailDoneMsg struct {
	note *apimodel.NoteModel
	err  error
}

func loginCmd(svc *appsvc.AppServices, user, pass string) tea.Cmd {
	return func() tea.Msg {
		return loginDoneMsg{err: svc.Login(context.Background(), user, pass)}
	}
}

func listCmd(svc *appsvc.AppServices, skip, take int) tea.Cmd {
	take = notes.NormalizeListTake(take)
	return func() tea.Msg {
		items, total, err := svc.ListNotes(context.Background(), skip, take)
		if err != nil {
			return listDoneMsg{err: err}
		}
		return listDoneMsg{items: items, total: total, skip: skip, take: take}
	}
}

func getNoteCmd(svc *appsvc.AppServices, noteID string) tea.Cmd {
	return func() tea.Msg {
		note, err := svc.GetNote(context.Background(), noteID)
		return noteDetailDoneMsg{note: note, err: err}
	}
}

func addCmd(svc *appsvc.AppServices, title, text string, tags []string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.AddNote(context.Background(), title, text, tags); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}

func updateCmd(svc *appsvc.AppServices, id, title, text string, tags []string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.UpdateNote(context.Background(), id, title, text, tags); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}

func deleteCmd(svc *appsvc.AppServices, id string) tea.Cmd {
	return func() tea.Msg {
		if err := svc.DeleteNote(context.Background(), id); err != nil {
			return simpleErrMsg{err: err}
		}
		return simpleOkMsg{}
	}
}
