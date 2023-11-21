# Octopoda简单时间协议 (Octopoda Simple Time Protocol)

用于保证run/xrun命令有相同的时间基准。

## 相关字段
1. NodeJoinResponse中的Ts: Brain的本地ms时间戳
2. HeartBeatRequest中的Delay: Tentacle估测的主从通信延迟
3. HeartBeatResponse中的Ts: Brain的本地ms时间戳
4. CommandParams中的ExecTs: Brain**预期tentacle同时执行命令的**的Brain本地ms时间戳

## 方法设计

### Step1 NodeJoin阶段
过程：
- Tentacle发送NodeJoinInfo，记录本地发送时间t1
- Brain收到NodeJoinInfo之后返回Brain本地时间戳tB
- Tentacle收到NodeJoinResponse，记录本地接收时间t2

Tentacle估测：
- Delay = (t2 - t1) / 2
- TsDiff = t1 + Delay - tB (表示自己和brain的本地时间戳的差)

### Step2 HeartBeat阶段
过程：
- Tentacle发送Delay，是之前估测的值
- Brain收到Delay，维护该Tentacle的延迟估计值
- 此过程保持Step1中的t1,t2,tB记录

### Step3 Command阶段
假设本次命令执行的目标节点为(N1,N2,...,Nx)，Brain维护的对应目标Delay为(Delay1,Delay2,...,Delayx)
Brain估测：
- ExecTs = Now() + lambda * max(Delay1,Delay2,...,Delayx)  (lambda >= 1, 视网络方差而定)
Tentacle估测：
- LocalExecTime = TsDiff + ExecTs
- LocalExecDefer = LocalExecTime - Now()

过程：
- Brain根据octl的指令确定目标节点集合，估测ExecTs并放入CommandParams中。
- Tentacle睡眠LocalExecDefer后开始执行命令。

