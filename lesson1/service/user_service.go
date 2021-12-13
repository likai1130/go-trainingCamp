package service

import (
	"github.com/pkg/errors"
	"go-trainingCamp/lesson1/dao"
	"go-trainingCamp/lesson1/entity"
	"log"
)

type UserService interface {
	InsertUserMany(users []entity.UserData) (count int, err error)
	FindUsers() (users []entity.UserData, err error)
	FindUser(filter map[string]interface{}) (user entity.UserData, err error)
}

var UserServiceInstance *userService

type userService struct {
	userDao dao.UserDao
}

func (u userService) InsertUserMany(users []entity.UserData) (count int, err error) {
	return u.userDao.InsertMany(users)
}

func (u userService) FindUsers() (users []entity.UserData, err error) {
	return u.userDao.FindAll()
}

func (u userService) FindUser(filter map[string]interface{}) (user entity.UserData, err error) {
	user, err = u.userDao.FindOne(filter)
	b, _ := dao.NewDbError().IsErrNoDocuments(err)
	if b {
		//数据不存在的处理
		var k string
		var v interface{}
		for key, value := range filter {
			k = key
			v = value
		}
		log.Printf("data not find. params is [%s = %v].Cause: %s\n", k, v, errors.Cause(err).Error())
		return user, nil
	}
	return user, err
}

func NewUserService() (userSvc UserService, err error) {
	if UserServiceInstance != nil {
		return UserServiceInstance, nil
	}
	userDaoInstance, err := dao.NewUserDao()
	if err != nil {
		return nil, err
	}
	UserServiceInstance = &userService{
		userDao: userDaoInstance,
	}
	return UserServiceInstance, nil
}
