package main

import (
	"fmt"
	"lesson3/bootstrap"
	"lesson3/config"
	"lesson3/middlewares/logging"
	mongodb "lesson3/pkg/mongdb"
	"lesson3/web/router"
)

func init() {
	config.InitConfig()
	mongodb.SetUp()
}

func newApp() *bootstrap.Bootstrapper {
	logging.SystemLogsSetUp()
	app := bootstrap.New("gin-web", "gin-web-example")
	app.Bootstrap()
	app.Configure(router.Configure)
	return app
}

//增加项目地址
// @termsOfService https://github.com/likai1130/go-example
func main() {
	app := newApp()
	conf := config.AppConfig
	port := conf.Server.Port
	listenPort := fmt.Sprintf(":%v", port) //格式化监听端口
	app.Listen(listenPort)

}
