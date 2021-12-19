package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/**
作业2： 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

实现方案:
1. 实现 HTTP server 的启动和关闭
2. 监听 linux signal信号，中断后级联退出goroutine,如Ctrl+C，结束程序，kill
3. errgroup 实现多个 goroutine 的级联退出

设计思路：
	使用errgroup管理多个goroutine，任意一个goroutine产生错误，则退出所有的goroutine。启动三个goroutine：
	go1：负责启动 http server
	go2：负责关闭 http server
	go3: 负责监听系统信号

总结：
	核心方法还是errgroup,errgroup包装了sync.WaitGroup,管理一组goroutine的声明周期。
	内部通过context.Background()创建了一个WithCancel，并将这个WithCancel的ctx返回，把Cancel封装起来用来取消上下文。
	内部使用sync.Once，一组goroutine的第一个error进行错误处理，只处理一次，即：调用ctx.CancelFunc
	group.Wait()等待所有goroutine结束后，再次调用g.cancel()并且返回第一个发生的错误。
	遗留问题：
		errgroup源码中为什么要做两次cancel，在group.Go()方法中的once方法中调用了一次g.cancel()，在group.wait()中调用了一次g.cancel()，
	对于同一个错误只cancel一次不就可以了么，我是用debug打断点的时候看到的。
*/
func main() {
	group, ctx := errgroup.WithContext(context.Background()) //创建errgroup，返回子ctx
	//创建http服务器并且监听8080端口
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello,world!"))
	})

	//模拟调用shutdown，级联退出g1,g2,g3
	shutdownSignal := make(chan struct{})
	mux.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
		shutdownSignal <- struct{}{}
	})
	server := &http.Server{Addr: ":8080", Handler: mux}

	//goroutine1 启动httpserver
	//模拟g1退出方法: 把Addr定义成字符串，g1就会启动http server失败。
	//g1退出后，group通过sync.Once会调用内部cancel方法使context不再阻塞
	//接着g2,g3退出
	group.Go(func() error {
		log.Println("http server start")
		return server.ListenAndServe()
	})

	//goroutine2
	//g2模拟错误方法：访问: http://127.0.0.1:8080/shutdown
	//g2 会调用Shutdown方法退出，因为server shutdown，所以接着g1退出；因为g2退出后，context不再阻塞，所以g3退出
	group.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exited1")
		case <-shutdownSignal:
			log.Println("http server shutdown ... ")
		}
		return server.Shutdown(ctx)
	})

	//goroutine3
	//g3 监控信号量，程序中断操作，所有goroutine结束 如：kill9，ctrl+C等所有相关信号
	//g3退出后，context不再阻塞了，g2会退出，退出时调用了shutdown，然后g1随之退出。
	group.Go(func() error {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		select {
		case <-ctx.Done():
			log.Println("errgroup exited2")
			return ctx.Err()
		case sig := <-signals:
			return errors.Errorf("get os signal: %v", sig)
		}
	})

	log.Printf("all go routine exited %+v\n", group.Wait())
}
