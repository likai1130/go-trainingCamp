package basic

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataSource interface {
	Insert(dataBaseName, collectName string, entityType interface{}) (string, error)
	InsertBatch(dataBaseName, collectName string, typeSlices []interface{}) (int, error)
	Get(dataBaseName, collectName string, filter bson.D) (*mongo.SingleResult, error)
	List(dataBaseName, collectName string, filter bson.D, options ...*options.FindOptions) ([]map[string]interface{}, error)
	Update(dataBaseName, collectName string, filter bson.D, entityType interface{}) (int64, error)
	UpdateBatch(dataBaseName, collectName string, filter bson.D, typeSlices interface{}) (int64, error)
	Delete(dataBaseName, collectName string, filter bson.D) (int64, error)
	DeleteBatch(dataBaseName, collectName string, filter bson.D) (int64, error)
	Aggregate(dataBaseName, collectName string, filter interface{}) ([]map[string]interface{}, error)
	Count(dataBaseName, collectName string, filter interface{}) (int64, error)
	FindOneAndUpdate(dataBaseName, collectName string, filter, update interface{}) (*mongo.SingleResult, error)
}
