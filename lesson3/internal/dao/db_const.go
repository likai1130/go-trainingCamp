package dao

import "time"

const (
	MONGODB_CONTEXT_TIMEOUT = 20 * time.Second //context超时时间
	MONGODB_DATABASE        = "camp"           //数据库名字
	MONGODB_COLLECT_USER    = "user"           //user表名

)
