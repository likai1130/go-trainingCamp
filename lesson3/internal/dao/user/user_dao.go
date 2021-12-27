package dao

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go-trainingCamp/lesson3/internal/dao"
	"go-trainingCamp/lesson3/internal/dao/basic"
	"go-trainingCamp/lesson3/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UserDao interface {
	InsertUser(user model.User) (string, error)
	GetUser(id string) (model.User, error)
	ListUser() ([]model.User, error)
	UpdateUser(data interface{}, id string) (int64, error)
	RemoveUser(id string) (int64, error)
}

type userDao struct {
	dataSource basic.DataSource
}

func (u *userDao) InsertUser(user model.User) (string, error) {
	return u.dataSource.Insert(dao.MONGODB_DATABASE, dao.MONGODB_COLLECT_USER, user)
}

func (u *userDao) GetUser(id string) (model.User, error) {
	filter := bson.D{{"_id", id}}
	result, err := u.dataSource.Get(dao.MONGODB_DATABASE, dao.MONGODB_COLLECT_USER, filter)
	if err != nil {
		return model.User{}, err
	}
	user := model.User{}
	err = result.Decode(&user)
	return user, err
}

func (u *userDao) ListUser() ([]model.User, error) {
	results, err := u.dataSource.List(dao.MONGODB_DATABASE, dao.MONGODB_COLLECT_USER, bson.D{})
	if err != nil {
		return nil, err
	}
	return resultsUserToUserSlices(results)
}

func resultsUserToUserSlices(results []map[string]interface{}) ([]model.User, error) {
	var owner []model.User
	if len(results) == 0 {
		return owner, nil
	}
	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, errors.Wrap(err, "resultsUserToUserSlices json.Marshal error")
	}
	if err = json.Unmarshal(marshal, &owner); err != nil {
		return nil, errors.Wrap(err, "resultsUserToUserSlices json.Unmarshal error")
	}
	return owner, err
}

func (u *userDao) UpdateUser(data interface{}, id string) (int64, error) {
	filter := bson.D{{"_id", id}}
	return u.dataSource.Update(dao.MONGODB_DATABASE, dao.MONGODB_COLLECT_USER, filter, data)
}

func (u *userDao) RemoveUser(id string) (int64, error) {
	filter := bson.D{{"_id", id}}
	return u.dataSource.Delete(dao.MONGODB_DATABASE, dao.MONGODB_COLLECT_USER, filter)
}


func NewUserDao() UserDao {
	return &userDao{
		dataSource: basic.NewSimpleImp(),
	}
}
