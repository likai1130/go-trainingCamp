
## Error Group 使用

作业2： 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

### 实现方案:

1. 实现 HTTP server 的启动和关闭
2. 监听 linux signal信号，中断后级联退出goroutine,如Ctrl+C，结束程序，kill
3. errgroup 实现多个 goroutine 的级联退出

### 设计思路：

使用errgroup管理多个goroutine，任意一个goroutine产生错误，则退出所有的goroutine。启动三个goroutine：
go1：负责启动 http server
go2：负责关闭 http server
go3: 负责监听系统信号

### 总结：

核心方法还是errgroup,errgroup包装了sync.WaitGroup,管理一组goroutine的声明周期。

内部通过context.Background()创建了一个WithCancel，并将这个WithCancel的ctx返回，把Cancel封装起来用来取消上下文。

内部使用sync.Once，一组goroutine的第一个error进行错误处理，只处理一次，即：调用ctx.CancelFunc
group.Wait()等待所有goroutine结束后，再次调用g.cancel()并且返回第一个发生的错误。

### 遗留问题：

errgroup源码中为什么要做两次cancel，在group.Go()方法中的once方法中调用了一次g.cancel()，在group.wait()中调用了一次g.cancel()，
对于同一个错误只cancel一次不就可以了么，我是用debug打断点的时候看到的。
