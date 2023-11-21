# 升级计划
1.6.8 -> 1.7.0

## 更新内容
主要内容为通信协议，包括：
1. httpns对外服务协议由https升级为http结合https，其中brain注册走https，其他走http。协议body格式也进行了更新。
2. octl与brain之间的通信协议由http升级为https，因为不兼容旧版的octl，octl必须升级
3. brain与tentacle之间的通信协议由tcp升级为tls，旧版不兼容
4. brain在tls通信中的角色发生了变化，由纯client变成名字服务的client和tentacle的server，产生的结果是需要采用新的配置
5. brain和httpns的证书由指定IP升级为指定名字，导致两者的证书必须更新

## 升级顺序
1. 更新brain和httpns证书
2. 在另一个端口启动新版httpNameServer(servePort=3555，DB=2)，保持旧版本运行。
3. 批量用sed命令修改pakma-tentacle的配置，将预览版时间由120s设置为600s，手动修改pakma-brain
4. 批量用sed命令修改tentacle的配置，包括：
   1. httpsNameServer.port=3555
5. 修改brain的配置，包括：
   1. httpsNameServer.port=3555
   2. clientCert: "/etc/octopoda/cert/client.pem" -> serverCert: "/etc/octopoda/cert/server.pem"
   3. clientKey: "/etc/octopoda/cert/client.key" -> serverKey: "/etc/octopoda/cert/server.key"
6. 为brain安装新的证书（纯服务端证书不再适用）：server.pem 和 server.key
7. 通过旧的Octopoda网络发送pakma升级命令
8. 修改本地octl的配置。包括：
   1. httpsNameServer.port=3555
9.  通过新的octl尝试get nodes，如果成功则confirm
10. httpNameServer旧版本停用，新版本重新部署到3455，tentacle/octl/brain配置文件改回3455
11. 3555端口的httpNameServer停用。

## 命令
```bash
# tentacle
cp /etc/octopoda/tentacle/tentacle.yaml /etc/octopoda/tentacle/tentacle_bak.yaml && sed -i 's/previewDuration: 120/previewDuration: 600/g' /etc/octopoda/tentacle/tentacle.yaml && sed -i 's/port: 3455/port: 3555/g' /etc/octopoda/tentacle/tentacle.yaml
```

## 异常预案
如果升级中途失败则考虑预案（所有tentacle新的配置中，名字服务的port发生了更改，变成了3555，其他可忽略）：
1. 手动将旧版本的httpns启动并监听3555端口
2. 手动回退brain和octl，确保octl get nodes能获取所有节点。
3. 排查故障并重试升级。
