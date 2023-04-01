package mongo

import "go.mongodb.org/mongo-driver/mongo"

func IsNoDocuments(err error) bool {
	return err == mongo.ErrNoDocuments
}
