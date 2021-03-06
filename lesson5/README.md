
## 第四次作业

### 1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

#### 测试参数

|参数选项|说明|
---|---|
|-h|指定服务器主机名。|
|-p|指定服务器端口。|
|-s|指定服务器 socket。|
|-c|指定并发连接数。|
|-n|指定请求的具体数量。。|
|-d|以字节的形式指定 SET/GET 值的数据大小。|
|-k|1 表示 keep alive；0 表示 reconnect，默认为 1|
|-r|SET/GET/INCR 使用随机 key, SADD 使用随机值。|
|-P|Pipeline 请求|
|-q|强制退出 Redis，仅显示 query/sec 值。|
|--csv|以 CSV 格式输出。|
|-l|生成循环，永久执行测试。|
|-t|仅运行以逗号分隔的测试命令列表。|
|-I(大写i)|空闲模式，打开 N 个空闲连接并等待连接。|

#### 测试命令

- redis-benchmark -t set,get -q -d 10
- redis-benchmark -t set,get -q -d 20
- redis-benchmark -t set,get -q -d 50
- redis-benchmark -t set,get -q -d 100
- redis-benchmark -t set,get -q -d 200
- redis-benchmark -t set,get -q -d 1k
- redis-benchmark -t set,get -q -d 5k

#### 测试结果

|方法|大小| 性能（次/秒） | 耗时(ms)
---|---|---|---
|SET|10B|18515.09|p50=2.215ms
|SET|20B|19312.48|p50=2.119ms
|SET|50B|19083.97|p50=2.159ms
|SET|100B|17724.21|p50=2.287ms
|SET|200B|19073.05|p50=2.143ms
|SET|1k|19087.61|p50=2.151ms
|SET|5k|14773.23|p50=2.711ms

|方法|大小| 性能（次/秒） | 耗时(ms)
---|---|---|---
|GET|10B|19853.09|p50=2.071ms
|GET|20B|19673.42|p50=2.095ms
|GET|50B|19047.62|p50=2.167ms
|GET|100B|18501.39|p50=2.207ms
|GET|200B|19065.78|p50=2.183ms
|GET|1k|19394.88|p50=2.119ms
|GET|5k|17053.21|p50=2.359ms

#### 2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

#### 测试命令
- redis-benchmark -t set -q -d 10 -n 100000
- redis-benchmark -t set -q -d 20 -n 100000
- redis-benchmark -t set -q -d 50 -n 100000
- redis-benchmark -t set -q -d 100 -n 100000
- redis-benchmark -t set -q -d 200 -n 100000
- redis-benchmark -t set -q -d 1k -n 100000
- redis-benchmark -t set -q -d 5k -n 100000
- ...

#### 测试SET写入的性能

|value大小| 写入次数 | Before Memory(kb)| After Memory(kb)| 单个key内存占用(B)
---|---|---|---|---
10B|10000|151008|150028|100.3
10B|100000|149756|149248|5.2
10B|200000|149992|147980|10.3
10B|300000|150016|149508|1.73
10B|400000|149524|149500|0.061
10B|500000|150220|149256|1.97


|value大小| 写入次数 | Before Memory(b)| After Memory(kb)| 单个key内存占用(B)
---|---|---|---|---
1k|10000|154140672|153636864|50.3
1k|100000|154124288|153128960|9.95
1k|200000|152850432|152330240|2.6
1k|300000|152834048|152838144|内存已被回收
1k|400000|153559040|152322048|3.0
1k|500000|152891392|149256|内存已经释放


|value大小| 写入次数 | Before Memory(b)| After Memory(b)| 单个key内存占用(b)
---|---|---|---|---
5k|10000|152350720|151629824|72
5k|100000|152092672|151072768|10.1
5k|200000|151080960|151080960|6.4
5k|300000|152367104|150839296|5.0
5k|400000|149524|149500|1.8
5k|500000|150220|149256|0.9


#### 总结

- 10b数据
  - 当写入数据为10b，次数1w次，平均key占用的大小为100b
  - 当写入数据为10b，次数40w次，平均key占用的大小为0.061b（内存占用最小）这个时候内存也在释放。
  - 对于10b数据整体来看，次数越大，平均key内存占用越少
- 1k数据
  - 当写入数据为1k，次数1w次，平均key占用的大小为50b，比10b数据少一倍。
  - 当写入数据为1k，次数20w次，平均key占用的大小为2.6（内存占用最小）
  - 因为内存边申请边释放，当50w的数据写入完毕后，内存就被释放掉了，这个很神奇，不知道为什么。
- 5k
  - 当写入数据为5k，次数1w次，平均key占用的大小为72b
  - 当写入数据为5k，次数50w次，平均key占用的大小为0.9b
  - 可见数据越大，次数越多，key占用的越小
- 总结
  - 综上所述，最佳数据值为5k的数据，数据量为50w,单个key占用内存最小




