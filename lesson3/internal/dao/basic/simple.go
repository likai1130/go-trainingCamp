package basic

import (
	"context"
	"github.com/pkg/errors"
	"go-trainingCamp/lesson3/common/logger"
	"go-trainingCamp/lesson3/internal/dao"
	mongodb "go-trainingCamp/lesson3/internal/pkg/mongdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var once sync.Once
var simpleInstance *simpleImp

type simpleImp struct {
	mongoClient *mongo.Client
}

func (s *simpleImp) Update(dataBaseName, collectName string, filter bson.D, entityType interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	result, err := collection.UpdateOne(ctx, filter, entityType)
	if err != nil {
		return 0, errors.Wrap(err, "update db data is error")
	}
	return result.ModifiedCount, nil
}

func (s *simpleImp) UpdateBatch(dataBaseName, collectName string, filter bson.D, typeSlices interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	results, err := collection.UpdateMany(ctx, filter, typeSlices)
	if err != nil {
		return 0, errors.Wrap(err, "update db many data is error ")
	}
	return results.ModifiedCount, nil
}

func (s *simpleImp) DeleteBatch(dataBaseName, collectName string, filter bson.D) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "delete db data is error ")
	}
	return result.DeletedCount, nil
}

func (s *simpleImp) Delete(dataBaseName, collectName string, filter bson.D) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "delete db data is error ")
	}
	return result.DeletedCount, nil
}

func (s *simpleImp) Insert(dataBaseName, collectName string, entityType interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	insertOneResult, err := collection.InsertOne(ctx, entityType)
	if err != nil {
		return "", errors.Wrap(err, "insert db  data is error")
	}
	insertId := insertOneResult.InsertedID.(primitive.ObjectID).Hex()
	return insertId, nil
}

func (s *simpleImp) InsertBatch(dataBaseName, collectName string, typeSlices []interface{}) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	insertManyResult, err := collection.InsertMany(ctx, typeSlices)
	if err != nil {
		return 0, errors.Wrap(err, "insert db many data is error ")
	}
	return len(insertManyResult.InsertedIDs), nil
}

func (s *simpleImp) List(dataBaseName, collectName string, filter bson.D, options ...*options.FindOptions) ([]map[string]interface{}, error) {
	if filter == nil {
		filter = bson.D{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	cursors, err := collection.Find(ctx, filter, options...)
	cursors.RemainingBatchLength()
	if err != nil || cursors.Err() != nil {
		return nil, errors.Wrap(err, "find db  data is error")
	}
	defer func(cursors *mongo.Cursor, ctx context.Context) {
		err := cursors.Close(ctx)
		if err != nil {
			logger.GetLogger().Errorf("List context error: %+v", errors.Wrap(err, "simpleImp list context close error"))
		}
	}(cursors, ctx)
	typesSlice := make([]map[string]interface{}, 0)
	err = cursors.All(ctx, &typesSlice)
	return typesSlice, err
}

func (s *simpleImp) Get(dataBaseName, collectName string, filter bson.D) (*mongo.SingleResult, error) {
	if filter == nil {
		filter = bson.D{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := s.mongoClient.Database(dataBaseName).Collection(collectName)
	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, errors.Wrap(result.Err(), "find db data is error")
	}
	return result, nil
}

func (c *simpleImp) FindOneAndUpdate(dataBaseName, collectName string, filter, update interface{}) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()
	collection := c.mongoClient.Database(dataBaseName).Collection(collectName)
	result := collection.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		return nil, errors.Wrap(result.Err(), "insertInc db data is error")
	}
	return result, nil
}

func (c *simpleImp) Count(dataBaseName, collectName string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := c.mongoClient.Database(dataBaseName).Collection(collectName)
	result, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.Wrap(err, "count is error")
	}
	return result, nil
}

// Aggregate 万能的聚合查询/**
func (c *simpleImp) Aggregate(dataBaseName, collectName string, filter interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dao.MONGODB_CONTEXT_TIMEOUT)
	defer cancel()

	collection := c.mongoClient.Database(dataBaseName).Collection(collectName)
	cursors, err := collection.Aggregate(ctx, filter)
	if err != nil || cursors.Err() != nil {
		return nil, errors.Wrap(err, "aggregate db is error")
	}
	defer func(cursors *mongo.Cursor, ctx context.Context) {
		err := cursors.Close(ctx)
		if err != nil {
			logger.GetLogger().Errorf("%+v", errors.Wrap(err, "complexImp list context close error"))
		}
	}(cursors, ctx)
	typesSlice := make([]map[string]interface{}, 0)
	err = cursors.All(ctx, &typesSlice)
	return typesSlice, err
}

func NewSimpleImp() DataSource {
	once.Do(func() {
		mongoCli, _ := mongodb.NewMongoCliInstance()
		simpleInstance = &simpleImp{
			mongoCli,
		}
	})
	return simpleInstance
}
