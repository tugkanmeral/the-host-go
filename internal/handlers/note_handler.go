package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	model "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/service"
	"github.com/tugkanmeral/the-host-go/internal/utils"
)

type NoteHandler struct {
	notes *service.NoteService
}

func NewNoteHandler(notes *service.NoteService) *NoteHandler {
	return &NoteHandler{notes: notes}
}

func (h *NoteHandler) Add(c *fiber.Ctx) error {
	var req model.NewNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Invalid body",
		})
	}

	if req.Title == "" || req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Title and Text fields are required",
		})
	}

	userID := utils.ExtractUserId(c)
	err := h.notes.Add(c.UserContext(), userID, req.Title, req.Text, req.Tags)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyOwnerID):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "UserId could not be extracted from token!",
			})
		case errors.Is(err, service.ErrNoteValidation):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "Title and Text fields are required",
			})
		case errors.Is(err, service.ErrNoteSaveFailed):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "Note could not be saved!",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "Note could not be saved!",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: "Note saved successfully",
	})
}

func (h *NoteHandler) GetList(c *fiber.Ctx) error {
	userID := utils.ExtractUserId(c)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	skipStr := c.Query("skip", "0")
	takeStr := c.Query("take", "10")

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		skip = 0
	}

	take, err := strconv.ParseInt(takeStr, 10, 64)
	if err != nil {
		take = 10
	}

	result, err := h.notes.GetList(c.UserContext(), userID, skip, take)
	if err != nil {
		if errors.Is(err, service.ErrEmptyOwnerID) {
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "UserId could not be extracted from token!",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to fetch notes!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.ListResponse[[]model.NoteListingItemModel]{
		Skip:       result.Skip,
		Take:       result.Take,
		TotalCount: result.TotalCount,
		Data:       result.Items,
	})
}

func (h *NoteHandler) Get(c *fiber.Ctx) error {
	userID := utils.ExtractUserId(c)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	id := c.Params("id", "")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Note ID is required!",
		})
	}

	noteModel, err := h.notes.Get(c.UserContext(), userID, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Message: "Note could not be found!",
			})
		case errors.Is(err, service.ErrNoteValidation):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "Note ID is required!",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "Note could not be found!",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(noteModel)
}

func (h *NoteHandler) Update(c *fiber.Ctx) error {
	userID := utils.ExtractUserId(c)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	id := c.Params("id", "")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Note ID is required!",
		})
	}

	var req model.NotePartialUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Invalid body",
		})
	}

	err := h.notes.Update(c.UserContext(), userID, id, req.Title, req.Text, req.Tags)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNoUpdateFields):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "At least one field (title, text, or tags) is required to update!",
			})
		case errors.Is(err, service.ErrNoteNotFound):
			return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
				Message: "Note could not be found!",
			})
		case errors.Is(err, service.ErrNoteValidation):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "Note ID is required!",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "Failed to update note!",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: "Note updated successfully",
	})
}

func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	userID := utils.ExtractUserId(c)
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	id := c.Params("id", "")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Note ID is required!",
		})
	}

	err := h.notes.Delete(c.UserContext(), userID, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to delete note!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Message: "Note deleted successfully",
	})
}
