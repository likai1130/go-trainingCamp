package api

import (
	"github.com/gin-gonic/gin"
	"go-trainingCamp/lesson3/common/app"
	"go-trainingCamp/lesson3/common/e"
	"go-trainingCamp/lesson3/internal/service"
	"go-trainingCamp/lesson3/internal/vo"
)

// User
// @Tags user
// @Summary 用户列表
// @Description 用户列表
// @Accept  json
// @Param lang header string false "国际化字段"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/users [GET]
func ListUsers(c *gin.Context) {
	app.NewResponse(service.NewUserService().ListUser(), c)
}

// User
// @Tags user
// @Summary 查询用户
// @Description 查询用户
// @Accept  json
// @Param lang header string false "国际化字段"
// @Param id path string true "用户ID"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/user/{id} [GET]
func GetUserById(c *gin.Context) {
	id := c.Param("id")
	app.NewResponse(service.NewUserService().GetUserById(id), c)
}

// User
// @Tags user
// @Summary 新增用户
// @Description 新增用户
// @Accept  json
// @Param lang header string false "国际化字段"
// @Param userParams body vo.UserParams true "用户对象"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/user [POST]
func SaveUser(c *gin.Context) {
	userParams := &vo.UserParams{}
	if err := c.ShouldBindJSON(userParams); err != nil {
		app.NewResponse(app.ResponseI18nMsgParams{Code: e.InvalidParams, Err: err}, c)
	} else {
		app.NewResponse(service.NewUserService().SaveUser(*userParams), c)
	}
}

// User
// @Tags user
// @Summary 删除用户
// @Description 删除用户
// @Accept  json
// @Param lang header string false "国际化字段"
// @Param id path string true "用户ID"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/user/{id} [DELETE]
func RemoveUser(c *gin.Context) {
	id := c.Param("id")
	app.NewResponse(service.NewUserService().RemoveUserById(id), c)
}

// User
// @Tags user
// @Summary 更新用户
// @Description 更新用户
// @Accept  json
// @Param lang header string false "国际化字段"
// @Param id path string true "用户ID"
// @Param userParams body vo.UserParams true "用户对象"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/user/{id} [PUT]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userParams := &vo.UserParams{}
	if err := c.ShouldBindJSON(userParams); err != nil {
		app.NewResponse(app.ResponseI18nMsgParams{Code: e.InvalidParams, Err: err}, c)
	} else {
		app.NewResponse(service.NewUserService().UpdateUserById(id, *userParams), c)
	}
}
