package mongodb

import (
	"context"
	"github.com/pkg/errors"
	"go-trainingCamp/lesson3/common/logger"
	"go-trainingCamp/lesson3/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

const MONGODB_CONNECT_TIMEOUT = 20 * time.Second //连接超时时间

var mongoCli *mongo.Client
var once sync.Once

func NewMongoCliInstance() (mongoInstance *mongo.Client, err error) {
	once.Do(func() {
		mongoCli, err = setUp()
	})
	return mongoCli, err
}

func setUp() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mConf := config.AppConfig.MongoConf
	clientOptions := options.Client().SetHosts(mConf.Hosts).
		SetMaxPoolSize(mConf.MaxPoolSize).         //最大连接数量
		SetConnectTimeout(MONGODB_CONNECT_TIMEOUT) //连接超时20s
	if mConf.UserName != "" && mConf.Password != "" {
		clientOptions.SetAuth(options.Credential{Username: mConf.UserName, Password: mConf.Password})
	}

	mongoInstance, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.WithMessage(err, "MongoDB connect fail")
	}

	err = mongoInstance.Ping(ctx, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "MongoDB ping fail")
	}
	logger.GetLogger().Info("MongoDB connect success !")
	return mongoInstance, nil
}

func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoCli.Disconnect(ctx)
}

func SetUp() {
	if _, err := NewMongoCliInstance(); err != nil {
		log.Fatal(err)
	}
}