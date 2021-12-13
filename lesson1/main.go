package main

import (
	"encoding/json"
	"go-trainingCamp/lesson1/entity"
	"go-trainingCamp/lesson1/service"
	"log"
)

var userSvc service.UserService

func init() {
	serviceInstance, err := service.NewUserService()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	userSvc = serviceInstance

}

func initUsersData() []entity.UserData {

	return []entity.UserData{
		entity.UserData{
			Name:       "kli1",
			Number:     152664,
			Age:        17,
			BirthMonth: 1852,
		},
		entity.UserData{
			Name:       "kli2",
			Number:     152665,
			Age:        18,
			BirthMonth: 1853,
		},
		entity.UserData{
			Name:       "kli3",
			Number:     152666,
			Age:        19,
			BirthMonth: 1854,
		},
	}

}

/**
	联合测试(模拟业务)
当这个调用方为上层服务

1. 插入3条数据
2. 查询全部
3. 查询一条不存在的数据
*/
func integrationTest() {
	usersData := initUsersData()

	//1. 插入数据
	count, err := userSvc.InsertUserMany(usersData)
	if err != nil {
		//记录日志
		log.Printf("第一步，批量插入数据失败: %+v \n", err)
		return
	}
	log.Printf("第一步，插入%d条数据成功\n", count)

	//2. 查询所有数据
	users, err := userSvc.FindUsers()
	if err != nil {
		//记录日志
		log.Printf("第二步，查询所有数据失败: %+v \n", err)
		return
	}
	marshal, _ := json.Marshal(users)
	log.Printf("第二步，查询所有user数据成功。users = %s \n", string(marshal))

	//3. 查询某个数据
	filter := map[string]interface{}{}
	filter["name"] = "zhangsan"
	user, err := userSvc.FindUser(filter)

	if err != nil {
		log.Printf("第三步，查询zhangsan失败: %s \n", err.Error())
		return
	}
	marshal, _ = json.Marshal(user)
	log.Printf("第三步，查询结果：user = %s \n", string(marshal))
	return
}

func main() {
	integrationTest()
}
