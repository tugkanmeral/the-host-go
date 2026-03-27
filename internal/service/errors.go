package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternal           = errors.New("internal error")
	ErrTokenGeneration    = errors.New("token generation failed")

	ErrUsernameTaken  = errors.New("username taken")
	ErrRegisterFailed = errors.New("register failed")

	ErrEmptyOwnerID       = errors.New("owner id required")
	ErrNoteNotFound       = errors.New("note not found")
	ErrNoteValidation     = errors.New("validation error")
	ErrNoteNoUpdateFields = errors.New("no fields to update")
	ErrNoteSaveFailed   = errors.New("note save failed")
	ErrNoteDeleteFailed = errors.New("note delete failed")
)
