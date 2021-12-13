package dao

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type DbError interface {
	IsErrNoDocuments(err error) (bool, error)
}

type dbError struct{}

func (d *dbError) IsErrNoDocuments(err error) (bool, error) {
	if errors.Cause(err) == mongo.ErrNoDocuments {
		return true, err
	}
	return false, err
}

func NewDbError() DbError {
	return &dbError{}
}
