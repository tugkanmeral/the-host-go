package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tugkanmeral/the-host-go/internal/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ExtractUserId(c *fiber.Ctx) string {
	authHeaderStr := c.Get("Authorization")
	jwt := strings.TrimPrefix(authHeaderStr, "Bearer ")
	userId := auth.GetUserId(jwt)
	return userId
}

func ConvertStringToObjectId(hex string) bson.ObjectID {
	objectID, err := bson.ObjectIDFromHex(hex)
	if err != nil {
		return bson.NilObjectID
	}
	return objectID
}

func ConvertObjectIdToString(objectId bson.ObjectID) string {
	return objectId.Hex()
}
