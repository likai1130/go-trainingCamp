package dto

import (
	"go-trainingCamp/lesson3/internal/biz/dict"
	"go-trainingCamp/lesson3/internal/model"
	"go-trainingCamp/lesson3/internal/vo"
)

//ModelUserToUserVO
func ModelUserToUserVO(user model.User) vo.UserVO {
	return vo.UserVO{
		Id:   user.Id,
		Name: user.Name,
		Age:  user.Age,
		Sex:  dict.GetSex(user.Sex),
	}
}

//ModelUserToUserVO
func ModelUserListToListUserVO(users []model.User) []vo.UserVO {
	var usersVO = []vo.UserVO{}
	for _, user := range users {
		userVO := vo.UserVO{
			Id:   user.Id,
			Name: user.Name,
			Age:  user.Age,
			Sex:  dict.GetSex(user.Sex),
		}
		usersVO = append(usersVO, userVO)
	}
	return usersVO
}