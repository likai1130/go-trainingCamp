package mclient

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const MONGODB_CONNECT_TIMEOUT  = 10 * time.Second

var cli *mongo.Client
var once sync.Once
func NewMongoClient() (mongoCli *mongo.Client, err error) {
	if cli == nil {
		client, err := setUp()
		if err != nil {
			return nil, errors.Wrap(err,"NewMongoClient error")
		}
		cli = client
	}
	return cli, nil
}

func setUp() (mongoCli *mongo.Client, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), MONGODB_CONNECT_TIMEOUT)
	defer cancelFunc()

	clientOptions := options.Client().SetHosts([]string{"127.0.0.1:27017"}).
		SetMaxPoolSize(100).         //最大连接数量
		SetConnectTimeout(MONGODB_CONNECT_TIMEOUT). //连接超时20s
		SetAuth(options.Credential{Username: "mongoadmin", Password: "secret"})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { //在正式项目中，这个错误不返回，应该直接失败
		return nil, errors.WithMessage(err, "connect mongo db fail")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.WithMessage(err, "ping mongo db fail")
	}

	return client, err
}