package service

import (
	"go-trainingCamp/lesson1/dao"
	"go-trainingCamp/lesson1/entity"
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
	return u.userDao.FindOne(filter)
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
