##  异常处理

**我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？**

答： 对于这个问题我的答案是这样的：sql.ErrNoRows常见的情况是根据条件Query某条数据，或者多条数据。 dao 层中当遇到一个 sql.ErrNoRows 的时候,不应该 Wrap 这个 error.而是在DAO中判断error类型是不是ErrNoRows，以错误降级的方式，用日志记录直接处理掉。

**我是这么理解的：**

1. **对于数据库操作分为两大种类型的错误：**
   
    数据库连接产生的错误。如：端口号问题，数据库地址连接不上，账号密码不对等，这些毫无疑问需要Wrap error，往上抛，甚至是在项目启动初始化时候就直接Fail掉。
    
    CRUD错误。对于CRUD的错误，如：主键不能为nil，插入的数据slice不能为nil，sql找不到数据等。有的需要向上抛，让调用方知道是产生错误的原因，具体处不处理，如何处理，看具体业务。但是查询不到数据的错误，在dao中可以屏蔽掉，不然调用方很难知道错误类型，如果逐个按需求判断会导致db中的错误类型引用的到处都是。


2. 我使用的是mongodb，所以用的是官方mongdb-driver，官方源码中mongo.errors包定义了CRUD返回的错误。

    ```
   // ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return
   // any documents.
   var ErrNoDocuments = errors.New("mongo: no documents in result")
   
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

    在mongodb中ErrNoDocuments错误类型和 sql.ErrNoRows是一样的，我就以mongodb为例子了。在mongdb-driver源码中，single_result.go文件中定义了ErrNoDocuments类型，这是一个全局的错误类型，可以在DAO中返回的错误类型中进行断言，然后打印日志做降级处理。
      
   ```
   var ErrNoDocuments = errors.New("mongo: no documents in result")
   ```
   我demo片段代码

   ```go
      func (u *userDao) FindOne(params bson.M)(userData entity.UserData, err error){
         ctx, cancelFunc := context.WithTimeout(context.Background(), MONGODB_CONNECT_TIMEOUT)
          defer cancelFunc()
      
          collection := u.cli.Database(MONGODB_DATABASE).Collection(MONGODB_DATABASE_COLLECT)
          result := collection.FindOne(ctx, params)
          if result.Err() != nil {
              //"数据找不到"的错误，降级
              if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
                  log.Println(errors.WithMessage(err,"dao not find data").Error())
                  return userData, nil
              }
              return userData,errors.Wrap(result.Err(),"find one error")
          }
      
          user := entity.UserData{}
          if err = result.Decode(&user); err != nil {
              return userData,errors.Wrap(err,"find one result.Decode error")
          }
          return user,err
      }
   ```
  
最后:

测试代码在main函数中，是main--service---dao---client的调用流程，模拟了简单业务查询操作。进入main.go进行查阅。