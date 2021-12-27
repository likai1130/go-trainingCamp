package dict

const (
	Male = 0
	Female = 1
)

var userSex = map[int]string{
	Male:"男",
	Female:"女",
}

func GetSex(sex int) string {
	return userSex[sex]
}