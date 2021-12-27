package service

import (
	"go-trainingCamp/lesson3/common/app"
	"go-trainingCamp/lesson3/common/logger"
	"go-trainingCamp/lesson3/internal/biz"
	"go-trainingCamp/lesson3/internal/vo"
)

type UserService interface {
	ListUser() app.ResponseI18nMsgParams
	GetUserById(id string) app.ResponseI18nMsgParams
	RemoveUserById(id string) app.ResponseI18nMsgParams
	UpdateUserById(id string, params vo.UserParams)app.ResponseI18nMsgParams
	SaveUser(params vo.UserParams) app.ResponseI18nMsgParams
}

var userInstacne  *userService

type userService struct {}

//ListUser 获取用户列表
func (u *userService) ListUser() app.ResponseI18nMsgParams {
	logger.GetLogger().Debug("请求用户列表")
	defer logger.GetLogger().Debug("请求用户列表结束")
	code, users, err := biz.NewUserBiz().ListUser()
	if err != nil {
		logger.GetLogger().Errorf("请求用户列表失败：err = %+v",err)
	}
	return app.NewI18nResponse(code, users, err)
}

//GetUserById 用户详情
func (u *userService) GetUserById(id string) app.ResponseI18nMsgParams {
	logger.GetLogger().Debugf("开始查询id=%s用户信息",id)
	defer logger.GetLogger().Debugf("结束查询id=%s用户信息",id)
	code, userVO, err := biz.NewUserBiz().GetUserById(id)
	if err != nil {
		logger.GetLogger().Errorf("获取用户 id = %s 信息失败：err = %+v",id,err)
	}
	return app.NewI18nResponse(code, userVO, err)
}

//RemoveUserById 删除用户
func (u *userService) RemoveUserById(id string) app.ResponseI18nMsgParams {
	logger.GetLogger().Debugf("开始删除id=%s用户信息",id)
	defer logger.GetLogger().Debugf("结束删除id=%s用户信息",id)
	code, userVO, err := biz.NewUserBiz().DeleteUserById(id)
	if err != nil {
		logger.GetLogger().Errorf("删除用户 id = %s 信息失败：err = %+v",id,err)
	}
	return app.NewI18nResponse(code, userVO, err)
}

//UpdateUserById 更新用户
func (u *userService) UpdateUserById(id string, params vo.UserParams) app.ResponseI18nMsgParams {
	logger.GetLogger().Debugf("开始更新id=%s用户信息",id)
	defer logger.GetLogger().Debugf("结束更新id=%s用户信息",id)
	code, userVO, err := biz.NewUserBiz().GetUserById(id)
	if err != nil {
		logger.GetLogger().Errorf("更新用户 id = %s 信息失败：err = %+v",id,err)
	}
	return app.NewI18nResponse(code, userVO, err)
}

//SaveUser新增用户
func (u *userService) SaveUser(params vo.UserParams) app.ResponseI18nMsgParams {
	logger.GetLogger().Debugf("开始新增用户 user = %+v",params)
	defer logger.GetLogger().Debugf("结束新增用户 user = %+v ",params)
	code, userVO, err := biz.NewUserBiz().SaveUser(params)
	if err != nil {
		logger.GetLogger().Errorf("新增用户失败：user = %+v, praerr = %+v",params,err)
	}
	return app.NewI18nResponse(code, userVO, err)
}

func NewUserService() UserService {
	if userInstacne == nil {
		userInstacne = &userService{}
	}
	return userInstacne
}