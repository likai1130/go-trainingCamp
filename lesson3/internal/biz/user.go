package biz

import (
	"go-trainingCamp/lesson3/common/e"
	"go-trainingCamp/lesson3/internal/biz/dto"
	"go-trainingCamp/lesson3/internal/errors"
	"go-trainingCamp/lesson3/internal/model"
	"go-trainingCamp/lesson3/internal/vo"
)

//ListUser
func (u *UserBiz) ListUser() (string, []vo.UserVO, error) {
	users, err := u.userDao.ListUser()
	if err != nil {
		return e.DBError, nil, err
	}
	return e.ScodeOK, dto.ModelUserListToListUserVO(users), nil
}

//GetUserById
func (u *UserBiz) GetUserById(id string) (string, vo.UserVO, error) {
	user, err := u.userDao.GetUser(id)
	if err != nil {
		if errors.IsErrNoDocuments(err) {
			return e.NotDocument, vo.UserVO{}, nil
		}
		return e.DBError, vo.UserVO{}, err
	}
	return e.ScodeOK, dto.ModelUserToUserVO(user), nil
}

//UpdateUserById
func (u *UserBiz) UpdateUserById(id string, user vo.UserParams) (string, int64, error) {
	updateUser, err := u.userDao.UpdateUser(user, id)
	if err != nil {
		return e.DBError, 0, err
	}
	return e.ScodeOK, updateUser, nil
}

//DeleteUserById
func (u *UserBiz) DeleteUserById(id string) (string, int64, error) {
	removeCount, err := u.userDao.RemoveUser(id)
	if err != nil {
		return e.DBError, 0, err
	}
	return e.ScodeOK, removeCount, err
}

//SaveUser
func (u *UserBiz) SaveUser(userParams vo.UserParams) (string, string, error) {
	user := model.User{
		Name: userParams.Name,
		Age:  userParams.Age,
		Sex:  userParams.Sex,
	}
	insertUser, err := u.userDao.InsertUser(user)
	if err != nil {
		return e.DBError, "", err
	}
	return e.ScodeOK, insertUser, nil
}
