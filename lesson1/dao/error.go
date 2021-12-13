package dao

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsErrNoDocuments(err error) bool {
	if errors.Cause(err) == mongo.ErrNoDocuments {
		return true
	}
	return false
}
