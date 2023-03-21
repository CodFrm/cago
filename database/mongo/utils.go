package mongo

import "go.mongodb.org/mongo-driver/mongo"

func IsNilDocument(err error) bool {
	return err == mongo.ErrNilDocument
}
