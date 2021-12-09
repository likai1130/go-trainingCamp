##  异常处理

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

答：

    老师好，对于这个问题我的答案是这样的：

    dao 层中当遇到一个 sql.ErrNoRows 的时候,是否应该 Wrap 这个 error.

    因为sql.ErrNoRows的情况有很多，就一个mongdb sql查询来说。如db异常导致连接失败、数据不存在no exist等
    dao层调用的是sql client，返回的是根错误，对于sql.ErrNoRows根错误时候，不需要使用Wrap来包装，直接返回调用方。