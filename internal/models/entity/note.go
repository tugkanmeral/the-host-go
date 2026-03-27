package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Note struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string        `bson:"title" json:"title"`
	Text           string        `bson:"text" json:"text"`
	Tags           []string      `bson:"tags" json:"tags"`
	CreationDate   time.Time     `bson:"creationDate" json:"creationDate"`
	LastUpdateDate time.Time     `bson:"lastUpdateDate" json:"lastUpdateDate"`
	OwnerId        string        `bson:"ownerId" json:"ownerId"`
}
