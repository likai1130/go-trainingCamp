package biz

import dao "go-trainingCamp/lesson3/internal/dao/user"

type UserBiz struct {
	userDao dao.UserDao
}

func NewUserBiz() *UserBiz {
	return &UserBiz{
		userDao: dao.NewUserDao(),
	}
}