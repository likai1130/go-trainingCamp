package model

type User struct {
	Id   string `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Age  int    `bson:"age" json:"age"`
	Sex  int    `bson:"sex" json:"sex"`
}
