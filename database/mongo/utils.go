package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

func IsNoDocuments(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}
