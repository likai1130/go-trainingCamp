package dao

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	mclient "go-trainingCamp/lesson1/client"
	"go-trainingCamp/lesson1/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	MONGODB_DATABASE         = "trainingcamp"
	MONGODB_DATABASE_COLLECT = "users"
	MONGODB_CONNECT_TIMEOUT  = 10 * time.Second
)

type UserDao interface {
	InsertMany(users []entity.UserData) (count int, err error)
	FindAll() (users []entity.UserData, err error)
	FindOne(filter bson.M) (user entity.UserData, err error)
}

var userDaoInstance *userDao

type userDao struct {
	cli *mongo.Client
}

func NewUserDao() (dao UserDao, err error) {
	if userDaoInstance != nil {
		return userDaoInstance, nil
	}
	client, err := mclient.NewMongoClient()
	if err != nil {
		return userDaoInstance, errors.WithMessage(err, "NewUserDao err")
	}
	userDaoInstance = &userDao{
		cli: client,
	}
	return userDaoInstance, nil
}

/**
插入
*/
func (u *userDao) InsertMany(userDatas []entity.UserData) (count int, err error) {
	var list []interface{}
	marshal, _ := json.Marshal(userDatas)
	json.Unmarshal(marshal, &list)

	ctx, cancelFunc := context.WithTimeout(context.Background(), MONGODB_CONNECT_TIMEOUT)
	defer cancelFunc()
	collection := u.cli.Database(MONGODB_DATABASE).Collection(MONGODB_DATABASE_COLLECT)
	results, err := collection.InsertMany(ctx, list)
	if err != nil {
		return 0, errors.Wrap(err, "insert many")
	}
	return len(results.InsertedIDs), err
}

/**
查询
*/
func (u *userDao) FindAll() (userData []entity.UserData, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), MONGODB_CONNECT_TIMEOUT)
	defer cancelFunc()
	collection := u.cli.Database(MONGODB_DATABASE).Collection(MONGODB_DATABASE_COLLECT)

	results, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return userData, errors.Wrap(err, "findAll error")
	}
	if results.Err() != nil {
		return userData, errors.Wrap(results.Err(), "findAll error")
	}

	data := []entity.UserData{}
	err = results.All(ctx, &data)
	if err != nil {
		return userData, errors.WithMessage(err, "findAll results.All error")
	}
	results.Close(ctx)
	return data, err
}

func (u *userDao) FindOne(params bson.M) (userData entity.UserData, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), MONGODB_CONNECT_TIMEOUT)
	defer cancelFunc()

	collection := u.cli.Database(MONGODB_DATABASE).Collection(MONGODB_DATABASE_COLLECT)
	result := collection.FindOne(ctx, params)
	if result.Err() != nil {
		return userData, errors.Wrap(result.Err(), "find one error")
	}

	user := entity.UserData{}
	if err = result.Decode(&user); err != nil {
		return userData, errors.Wrap(err, "find one result.Decode error")
	}
	return user, err
}
