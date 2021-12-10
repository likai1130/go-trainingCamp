##  异常处理

**我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？**

答： 对于这个问题我的答案是这样的： dao 层中当遇到一个 sql.ErrNoRows 的时候,应该 Wrap 这个 error.

**我是这么理解的：**

1. **对于数据库操作分为两大种类型的错误：**
   
    数据库连接产生的错误。如：端口号问题，数据库地址连接不上，账号密码不对等，这些毫无疑问需要Wrap error，往上抛，甚至是在项目启动初始化时候就直接Fail掉。
    
    CRUD错误。对于CRUD的错误，如：主键不能为nil，插入的数据slice不能为nil，sql找不到数据等。也是需要向上抛，让调用方知道是产生错误的原因，具体处不处理，如何处理，看具体业务。

2. 我使用的是mongodb，所以用的是官方mongdb-driver，官方源码中mongo.errors包定义了CRUD返回的错误。

    ```
    // ErrNilDocument is returned when a nil document is passed to a CRUD method.
    var ErrNilDocument = errors.New("document is nil")
    
    // ErrNilValue is returned when a nil value is passed to a CRUD method.
    var ErrNilValue = errors.New("value is nil")
    
    // ErrEmptySlice is returned when an empty slice is passed to a CRUD method that requires a non-empty slice.
    var ErrEmptySlice = errors.New("must provide at least one element in input slice")
    
    // ErrMapForOrderedArgument is returned when a map with multiple keys is passed to a CRUD method for an ordered parameter
    type ErrMapForOrderedArgument struct {
    ParamName string
    }
    
    ```

    可以看出来对于CRUD返回的sql.ErrNoRows的情况有很多，基本上都是操作不当，数据有问题，以及查询不到结果等错误，这意味着可以在DAO层将err和sql.ErrNoRows都Wrap，然后向上抛，在调用方根据error类做错误降级的处理。