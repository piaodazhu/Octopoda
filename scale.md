## 关于扩展
### httpsNameServer
在架构中，httpsNameServer面临较大的压力，原因在于:
1. https握手和加密开销较大。
2. 截至目前(v1.5.2版本)，tentacle需要周期性地向httpsNameServer请求Token。

要支持大规模(10k)的网络，httpsNameServer可以如下扩展：
1. 考虑到安全性，https保持不变。但是采用长连接实现，减少握手次数。
2. httpsNameServer读写分离，主节点负责名字注册和Token生产，从节点从主节点定时同步Token，tentacle从从节点同步Token。
3. httpsNameServer的redis后端也设置为对应的主从结构。

httpsNameServer依赖少，逻辑简单，挂掉的几率很低。但是一个可以优化的点是：通过合理设计tentacle和brain的逻辑，确保httpsNameServer即使挂掉了，在一定时间(如10min)内恢复，网络不受任何影响。即，httpsNameServer挂掉了，Token就从此不刷新，brain的名字也从此不更新，brain和tentacle哪怕连不通httpsNameServer，他们之间也继续用故障前的信息保持通信。

### Brain
在架构中，Brain日常面临的压力主要在于要跟所有tentacle建立2条长连接，通信也需要用Token来加密解密。至于一个Brain能承受多少个tentacle，这个是没有测过，不确定的。从Goroutine并发能力的角度，10k应该轻轻松松，但是实际加上加密解密、序列化反序列化、逻辑处理、日志打印等，可能比较吃力。目前的想法是:
1. 针对日常的心跳包优先做优化，减少该链路的序列化、日志等处理，尽可能纯TCP简单处理。
2. 针对model里面记录所有节点大表进行优化，用哈希到多锁、或者sync.Map来优化。

Brain可能挂掉，这对网络的影响是很严重的。对此可以采用热备的方式：
1. N个Brain运行，共用一个name，为的是确保不同的Brain在其他模块眼里是同一个。N个Brain有各自的UUID，为的是唯一标识自己在Brain之间的身份。
2. 同时只有1个Brain真正在工作。工作的Brain周期性向httpsNameServer注册名字(目前已有的逻辑)，同时做一个分布式锁续约的操作(需要实现的，利用redis后端，用自己的UUID做value)。
3. 其他热备Brain周期性地抢锁，抢到了就变成新的工作Brain，注册自己的名字并续约。
4. 设置：续约为1s，分布式锁超时时间2s，抢锁为0.1s。总体下来，有Brain挂掉，最多2s后就能自动切换。
5. 问题1：如果Brain下线，却和tentacle继续保持连接，那么tentacle是不是不知道Brain已经变成新的了？可以在Token中携带一个更新Brain的信号，强制tentacle重新解析Brain的名字。
6. 问题2：Brain之间的一致性。目前(v1.5.2版本)Brain需要持久化维护场景版本信息和组信息，组信息放redis里面，场景放在磁盘文件里。解决的方式是都存在同一个redis实例或者集群里。

### Tentacle
tentacle不面临扩展的压力，但是最好是在功能稳定后减少后台运行的开销。去除不必要的序列化和打印操作。
