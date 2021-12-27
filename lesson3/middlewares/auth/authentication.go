package auth

import (
	"github.com/gin-gonic/gin"
	"go-trainingCamp/lesson3/common/app"
	"go-trainingCamp/lesson3/common/e"
	"go-trainingCamp/lesson3/common/logger"
)

/**
 * http api 鉴权
 */
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		authorization := c.GetHeader("Authorization")
		logger.GetLogger().Infof("当前请求 authorization = %s", authorization)
		//TODO 数据库验证authorization是否存在

		if err != nil {
			app.NewResponse(app.ResponseI18nMsgParams{C: c, Code: e.StatusUnauthorized, Err: err})
			c.Abort()
			return
		}
		c.Next()
	}
}
