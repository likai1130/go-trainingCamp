package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-trainingCamp/lesson3/common/logger"
	"go-trainingCamp/lesson3/middlewares/crossdomain"
	"go-trainingCamp/lesson3/middlewares/logging"
	"go-trainingCamp/lesson3/middlewares/metrics"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*gin.Engine
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
}

// New returns a new Bootstrapper.
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Engine:       gin.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}
func (b *Bootstrapper) Bootstrap() *Bootstrapper {

	//设置业务日志级别或者中间件
	b.Use(logging.LoggerToFile(b.AppName))

	//设置监控
	if p := metrics.PrometheusSetUp(); p != nil {
		p.Use(b.Engine)
	}
	b.Use(gin.Recovery())
	b.Use(crossdomain.Cors())
	//b.Use(middlewares.Authentication())

	return b
}

//优雅终止
func (b *Bootstrapper) Listen(addr string, cfgs ...Configurator) {
	server := &http.Server{
		Addr:    addr,
		Handler: b,
	}
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		logger.GetLogger().Info("server start ...")
		return server.ListenAndServe()
	})

	group.Go(func() error {
		select {
		case <-ctx.Done():
			logger.GetLogger().Errorf("server shut down  ...")
			return server.Shutdown(ctx)
		}
	})

	group.Go(func() error {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case sig := <-signals:
				return errors.Errorf("get os signal: %v", sig)
			}
	})

	err := group.Wait()
	logger.GetLogger().Errorf("server down err: %v", err)

	//b.Run(addr)
}
