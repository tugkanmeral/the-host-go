package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tugkanmeral/the-host-go/internal/database"
	model "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/models/entity"
	"github.com/tugkanmeral/the-host-go/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func Add(c *fiber.Ctx) error {
	var req model.NewNoteRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Invalid body",
		})
	}

	if req.Title == "" || req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Title and Text fields are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userId := utils.ExtractUserId(c)
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	db := database.GetDB()
	newNote := &entity.Note{
		Title:        req.Title,
		Text:         req.Text,
		Tags:         req.Tags,
		OwnerId:      userId,
		CreationDate: time.Now(),
	}

	_, err := db.Collection(database.NoteCollectionName).InsertOne(ctx, newNote)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Note could not be saved!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Data: "Note saved successfully",
	})
}

func GetList(c *fiber.Ctx) error {
	userId := utils.ExtractUserId(c)
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "UserId could not be extracted from token!",
		})
	}

	skipStr := c.Params("skip", "0")
	takeStr := c.Params("take", "10")

	skip, err := strconv.ParseInt(skipStr, 10, 64)
	if err != nil {
		skip = 0
	}

	take, err := strconv.ParseInt(takeStr, 10, 64)
	if err != nil {
		take = 10
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := database.GetDB()

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"ownerId": userId,
			},
		},
		bson.M{
			"$sort": bson.M{
				"creationDate": -1,
			},
		},
		bson.M{
			"$skip": skip,
		},
		bson.M{
			"$limit": take,
		},
	}

	cursor, err := db.Collection(database.NoteCollectionName).Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to fetch notes!",
		})
	}
	defer cursor.Close(ctx)

	var notes []entity.Note
	if err = cursor.All(ctx, &notes); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to decode notes!",
		})
	}

	var noteListings []model.NoteListingItemModel
	for _, note := range notes {
		noteListings = append(noteListings, model.NoteListingItemModel{
			Id:    note.ID.Hex(),
			Title: note.Title,
			Text:  note.Text,
			Tags:  note.Tags,
		})
	}

	if noteListings == nil {
		noteListings = []model.NoteListingItemModel{}
	}

	totalCount, err := db.Collection(database.NoteCollectionName).CountDocuments(ctx, bson.M{
		"ownerId": userId,
	})
	if err != nil {
		fmt.Println(err.Error())
		totalCount = 0
	}

	return c.Status(fiber.StatusOK).JSON(model.ListResponse[[]model.NoteListingItemModel]{
		Skip:       int(skip),
		Take:       int(take),
		TotalCount: int(totalCount),
		Data:       noteListings,
	})
}

func Get(c *fiber.Ctx) error {
	userId := utils.ExtractUserId(c)
	if userId == "" {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId := utils.ConvertStringToObjectId(id)
	filter := bson.D{
		{Key: "ownerId", Value: userId},
		{Key: "_id", Value: objectId},
	}

	db := database.GetDB()
	var note entity.Note
	err := db.Collection(database.NoteCollectionName).FindOne(ctx, filter).Decode(&note)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Message: "Note could not be found!",
		})
	}

	noteModel := &model.NoteModel{
		Id:             utils.ConvertObjectIdToString(note.ID),
		Title:          note.Title,
		Text:           note.Text,
		Tags:           note.Tags,
		CreationDate:   note.CreationDate,
		LastUpdateDate: note.LastUpdateDate,
	}

	return c.Status(fiber.StatusOK).JSON(noteModel)
}

func Update(c *fiber.Ctx) error {
	userId := utils.ExtractUserId(c)
	if userId == "" {
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
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Invalid body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId := utils.ConvertStringToObjectId(id)
	filter := bson.D{
		{Key: "ownerId", Value: userId},
		{Key: "_id", Value: objectId},
	}

	// Build update document with only non-nil/non-empty fields
	updateFields := bson.M{}
	if req.Title != "" {
		updateFields["title"] = req.Title
	}
	if req.Text != "" {
		updateFields["text"] = req.Text
	}
	if len(req.Tags) > 0 {
		updateFields["tags"] = req.Tags
	}
	// Always update the lastUpdateDate
	updateFields["lastUpdateDate"] = time.Now()

	// If no fields to update, return error
	if len(updateFields) == 1 { // Only lastUpdateDate
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "At least one field (title, text, or tags) is required to update!",
		})
	}

	db := database.GetDB()
	result, err := db.Collection(database.NoteCollectionName).UpdateOne(
		ctx,
		filter,
		bson.M{"$set": updateFields},
	)

	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to update note!",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Message: "Note could not be found!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Data: "Note updated successfully",
	})
}

func Delete(c *fiber.Ctx) error {
	userId := utils.ExtractUserId(c)
	if userId == "" {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId := utils.ConvertStringToObjectId(id)
	filter := bson.D{
		{Key: "ownerId", Value: userId},
		{Key: "_id", Value: objectId},
	}

	db := database.GetDB()
	res, err := db.Collection(database.NoteCollectionName).DeleteOne(ctx, filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to delete note!",
		})
	}

	if res.DeletedCount < 1 {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Failed to delete note!",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SuccessResponse{
		Data: "Note deleted successfully",
	})
}
