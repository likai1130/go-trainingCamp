##  异常处理

**我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？**

答： 对于这个问题我的答案是这样的：sql.ErrNoRows常见的情况是根据条件Query某条数据，或者多条数据。 dao 层中当遇到一个 sql.ErrNoRows 的时候,应该 Wrap 这个 error.往上抛，通过接口方法对外暴露sql.ErrNoRows错误信息。

**我是这么理解的：**

1. **对于数据库操作分为两大种类型的错误：**
   
    数据库连接产生的错误。如：端口号问题，数据库地址连接不上，账号密码不对等，这些毫无疑问需要Wrap error，往上抛，甚至是在项目启动初始化时候就直接Fail掉。
    
    CRUD错误。对于CRUD的错误，如：主键不能为nil，插入的数据slice不能为nil，sql找不到数据等。有的需要向上抛，让调用方知道是产生错误的原因，具体处不处理，如何处理，看具体业务。但是查询不到数据的错误，需要在内部提供错误判断方法，是否处理由调用方决定。


2. 我使用的是mongodb，所以用的是官方mongdb-driver，官方源码中mongo.errors包定义了CRUD返回的错误。

    ```go
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

    在mongodb中ErrNoDocuments错误类型和 sql.ErrNoRows是一样的，我就以mongodb为例子了。在mongdb-driver源码中，single_result.go文件中定义了ErrNoDocuments类型，这是一个全局的错误类型，可以在DAO中进行错误类型中进行断言，将断言结果告诉调用方。。
   single_result.go文件：
   
   ```go
    // ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return
   // any documents.
    var ErrNoDocuments = errors.New("mongo: no documents in result")
   
   ```
   我demo片段代码
   ```go
   import (
       "github.com/pkg/errors"
       "go.mongodb.org/mongo-driver/mongo"
      )
   
      type DbError interface {
         IsErrNoDocuments(err error) (bool, error)
      }
      
      type dbError struct{}
      
      func (d *dbError) IsErrNoDocuments(err error) (bool, error) {
         if errors.Cause(err) == mongo.ErrNoDocuments {
         return true, err
      }
         return false, err
      }
      
      func NewDbError() DbError {
         return &dbError{}
      }
   ```
   service 

   ```go
   func (u userService) FindUser(filter map[string]interface{}) (user entity.UserData, err error) {
       user, err = u.userDao.FindOne(filter)
       b, _ := dao.NewDbError().IsErrNoDocuments(err)
       if b {
           //数据不存在的处理
           var k string
           var v interface{}
           for key, value := range filter {
               k = key
               v = value
           }
           log.Printf("data not find. params is [%s = %v].Cause: %s\n", k, v, errors.Cause(err).Error())
           return user, nil
       }
       return user, err
   }
   
   ```
  
最后:

测试代码在main函数中，是main--service---dao---client的调用流程，模拟了简单业务查询操作。进入main.go进行查阅。