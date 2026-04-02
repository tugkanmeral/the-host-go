package service

import (
	"context"
	"errors"
	"time"

	"github.com/tugkanmeral/the-host-go/internal/database"
	model "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/models/entity"
	"github.com/tugkanmeral/the-host-go/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NoteListResult struct {
	Items      []model.NoteListingItemModel
	Skip       int
	Take       int
	TotalCount int
}

type NoteService struct {
	db *mongo.Database
}

func NewNoteService(db *mongo.Database) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) Add(ctx context.Context, ownerID, title, text string, tags []string) error {
	if ownerID == "" {
		return ErrEmptyOwnerID
	}
	if title == "" || text == "" {
		return ErrNoteValidation
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	newNote := &entity.Note{
		Title:        title,
		Text:         text,
		Tags:         tags,
		OwnerId:      ownerID,
		CreationDate: time.Now(),
	}

	_, err := s.db.Collection(database.NoteCollectionName).InsertOne(ctx, newNote)
	if err != nil {
		return ErrNoteSaveFailed
	}

	return nil
}

func (s *NoteService) GetList(ctx context.Context, ownerID string, skip, take int64, searchTerm string) (*NoteListResult, error) {
	if ownerID == "" {
		return nil, ErrEmptyOwnerID
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var pipeline bson.A

	// 1. MATCH
	matchStage := bson.M{"ownerId": ownerID}
	if searchTerm != "" {
		matchStage["$text"] = bson.M{"$search": searchTerm}
	}
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// 2. SORT
	sortStage := bson.D{}
	if searchTerm != "" {
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"score": bson.M{"$meta": "textScore"},
			},
		})
		sortStage = bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}}
	} else {
		sortStage = bson.D{{Key: "creationDate", Value: -1}}
	}
	pipeline = append(pipeline, bson.M{"$sort": sortStage})

	// 3. PAGINATION
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": take})

	cursor, err := s.db.Collection(database.NoteCollectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, ErrInternal
	}
	defer cursor.Close(ctx)

	var notes []entity.Note
	if err = cursor.All(ctx, &notes); err != nil {
		return nil, ErrInternal
	}

	noteListings := make([]model.NoteListingItemModel, 0, len(notes))
	for _, note := range notes {
		text := note.Text
		if len(note.Text) > 25 {
			text = note.Text[:25] + "..."
		}
		noteListings = append(noteListings, model.NoteListingItemModel{
			Id:    note.ID.Hex(),
			Title: note.Title,
			Text:  text,
			Tags:  note.Tags,
		})
	}

	countFilter := bson.M{"ownerId": ownerID}
	if searchTerm != "" {
		countFilter["$text"] = bson.M{"$search": searchTerm}
	}
	totalCount, err := s.db.Collection(database.NoteCollectionName).CountDocuments(ctx, countFilter)
	if err != nil {
		totalCount = 0
	}

	return &NoteListResult{
		Items:      noteListings,
		Skip:       int(skip),
		Take:       int(take),
		TotalCount: int(totalCount),
	}, nil
}

func (s *NoteService) Get(ctx context.Context, ownerID, noteID string) (*model.NoteModel, error) {
	if ownerID == "" {
		return nil, ErrEmptyOwnerID
	}
	if noteID == "" {
		return nil, ErrNoteValidation
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID := utils.ConvertStringToObjectId(noteID)
	filter := bson.D{
		{Key: "ownerId", Value: ownerID},
		{Key: "_id", Value: objectID},
	}

	var note entity.Note
	err := s.db.Collection(database.NoteCollectionName).FindOne(ctx, filter).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoteNotFound
		}
		return nil, ErrInternal
	}

	return &model.NoteModel{
		Id:             utils.ConvertObjectIdToString(note.ID),
		Title:          note.Title,
		Text:           note.Text,
		Tags:           note.Tags,
		CreationDate:   note.CreationDate,
		LastUpdateDate: note.LastUpdateDate,
	}, nil
}

func (s *NoteService) Update(ctx context.Context, ownerID, noteID string, title, text string, tags []string) error {
	if ownerID == "" {
		return ErrEmptyOwnerID
	}
	if noteID == "" {
		return ErrNoteValidation
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID := utils.ConvertStringToObjectId(noteID)
	filter := bson.D{
		{Key: "ownerId", Value: ownerID},
		{Key: "_id", Value: objectID},
	}

	updateFields := bson.M{}
	if title != "" {
		updateFields["title"] = title
	}
	if text != "" {
		updateFields["text"] = text
	}
	if len(tags) > 0 {
		updateFields["tags"] = tags
	}
	updateFields["lastUpdateDate"] = time.Now()

	if len(updateFields) == 1 {
		return ErrNoteNoUpdateFields
	}

	result, err := s.db.Collection(database.NoteCollectionName).UpdateOne(
		ctx,
		filter,
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return ErrInternal
	}

	if result.MatchedCount == 0 {
		return ErrNoteNotFound
	}

	return nil
}

func (s *NoteService) Delete(ctx context.Context, ownerID, noteID string) error {
	if ownerID == "" {
		return ErrEmptyOwnerID
	}
	if noteID == "" {
		return ErrNoteValidation
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objectID := utils.ConvertStringToObjectId(noteID)
	filter := bson.D{
		{Key: "ownerId", Value: ownerID},
		{Key: "_id", Value: objectID},
	}

	res, err := s.db.Collection(database.NoteCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		return ErrNoteDeleteFailed
	}

	if res.DeletedCount < 1 {
		return ErrNoteDeleteFailed
	}

	return nil
}
