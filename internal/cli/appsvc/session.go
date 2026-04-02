package appsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/tugkanmeral/the-host-go/internal/auth"
	apimodel "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/service"
)

// AppServices wraps service layer access for the TUI; authenticated notes use userID from the last successful login.
type AppServices struct {
	Auth   *service.AuthService
	Notes  *service.NoteService
	userID string
}

func NewAppServices(authSvc *service.AuthService, noteSvc *service.NoteService) *AppServices {
	return &AppServices{Auth: authSvc, Notes: noteSvc}
}

func (a *AppServices) Login(ctx context.Context, username, password string) error {
	token, err := a.Auth.Login(ctx, username, password)
	if err != nil {
		return formatServiceError(err)
	}
	uid := auth.GetUserId(token)
	if uid == "" {
		return fmt.Errorf("could not resolve session user id")
	}
	a.userID = uid
	return nil
}

func (a *AppServices) ListNotes(ctx context.Context, skip, take int, searchTerm string) ([]apimodel.NoteListingItemModel, int, error) {
	result, err := a.Notes.GetList(ctx, a.userID, int64(skip), int64(take), searchTerm)
	if err != nil {
		return nil, 0, formatServiceError(err)
	}
	return result.Items, result.TotalCount, nil
}

func (a *AppServices) GetNote(ctx context.Context, noteID string) (*apimodel.NoteModel, error) {
	note, err := a.Notes.Get(ctx, a.userID, noteID)
	if err != nil {
		return nil, formatServiceError(err)
	}
	return note, nil
}

func (a *AppServices) AddNote(ctx context.Context, title, text string, tags []string) error {
	err := a.Notes.Add(ctx, a.userID, title, text, tags)
	if err != nil {
		return formatServiceError(err)
	}
	return nil
}

func (a *AppServices) UpdateNote(ctx context.Context, id, title, text string, tags []string) error {
	err := a.Notes.Update(ctx, a.userID, id, title, text, tags)
	if err != nil {
		return formatServiceError(err)
	}
	return nil
}

func (a *AppServices) DeleteNote(ctx context.Context, id string) error {
	err := a.Notes.Delete(ctx, a.userID, id)
	if err != nil {
		return formatServiceError(err)
	}
	return nil
}

func formatServiceError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return fmt.Errorf("invalid username or password")
	case errors.Is(err, service.ErrTokenGeneration):
		return fmt.Errorf("failed to generate token")
	case errors.Is(err, service.ErrInternal):
		return fmt.Errorf("internal server error")
	case errors.Is(err, service.ErrEmptyOwnerID):
		return fmt.Errorf("user id missing; please sign in again")
	case errors.Is(err, service.ErrNoteNotFound):
		return fmt.Errorf("note not found")
	case errors.Is(err, service.ErrNoteValidation):
		return fmt.Errorf("invalid note data")
	case errors.Is(err, service.ErrNoteNoUpdateFields):
		return fmt.Errorf("at least one field (title, text, or tags) is required to update")
	case errors.Is(err, service.ErrNoteSaveFailed):
		return fmt.Errorf("could not save note")
	case errors.Is(err, service.ErrNoteDeleteFailed):
		return fmt.Errorf("could not delete note or note was not found")
	default:
		return err
	}
}
