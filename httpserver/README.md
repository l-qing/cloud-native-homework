# 模块八作业
作业要求

编写 Kubernetes 部署脚本将 httpserver 部署到 kubernetes 集群，以下是你可以思考的维度

优雅启动
优雅终止
资源需求和 QoS 保证
探活
日常运维需求，日志等级
配置和代码分离
尝试用Service、Ingress将你的服务发布给集群外部的调用方
确保整个应用的高可用
通过证书保证httpServer的通讯安全
提交截止时间：11 月 28 日（下周日） 23:59

## 作业说明
* 优雅启动：模拟了一个启动时的耗时，便于测试启动探针等。也配置了亲和性属性，保证多副本不调度到一个节点上。
* 优雅终止：程序捕获了信号，并做了幽雅退出操作。Pod也调整了最大超时时间。还可以增加 preStop 钩子。
* 资源需求：配置了容器的 resources 信息。
* 探活：新增了启动探针、存活探针、就绪探针。
* 常运维需求，日志等级：暂时只支持指定参数的方式调整日志级别，暂未实现从configmap中读取该信息。
* 配置和代码分离：版本信息从env、secret、configmap中配置，并在程序中读取。模拟简单的分离。
* 高可用：设置了容器副本数，且配置了亲和性。避免因为单节点故障导致服务不可用。（暂时不通过dns实现跨集群的冗余部署。）
* 外部访问：用NodePord Service、及 Ingress 将服务发布给集群外部的调用方

TODO:
* 可以从configmap中读取日志配置信息
* 可以配置日志收集服务，将日志集中收集。
* 可以配置 HPA，在容器负载发生变化时实现动态扩缩容，减少运维工作量。

## 使用说明
1. 推送镜像

```shell
$ VERSION=2.0.1 make build push

docker build --build-arg version=v2.1.0 . -t 190219044/httpserver:2.0.1
docker push 190219044/httpserver:2.0.1
```

2. 创建 k8s 资源

```shell
$ kubectl apply -f deployment.yaml

deployment.apps/cloudnative-httpserver created
service/cloudnative-httpserver created
configmap/cloudnative-httpserver-configmap created
secret/cloudnative-httpserver-secret created

$ kubectl logs -f -l app=cloudnative-httpserver

```

3. 滚动更新

```shell
$ kubectl set image deployment/cloudnative-httpserver httpserver=190219044/httpserver:2.0.2

deployment.apps/cloudnative-httpserver image updated

$ kubectl get events -w

0s          Normal    ScalingReplicaSet         deployment/cloudnative-httpserver              Scaled up replica set cloudnative-httpserver-69fffbc5c7 to 1
0s          Normal    SuccessfulDelete          replicaset/cloudnative-httpserver-555bdf4f5d   Deleted pod: cloudnative-httpserver-555bdf4f5d-vnl29
0s          Normal    ScalingReplicaSet         deployment/cloudnative-httpserver              Scaled down replica set cloudnative-httpserver-555bdf4f5d to 0
0s          Normal    Killing                   pod/cloudnative-httpserver-555bdf4f5d-vnl29    Stopping container httpserver
0s          Normal    SuccessfulCreate          replicaset/cloudnative-httpserver-69fffbc5c7   Created pod: cloudnative-httpserver-69fffbc5c7-sgmmr
0s          Normal    Scheduled                 pod/cloudnative-httpserver-69fffbc5c7-sgmmr    Successfully assigned devops/cloudnative-httpserver-69fffbc5c7-sgmmr to 192.168.130.63
0s          Normal    Pulling                   pod/cloudnative-httpserver-69fffbc5c7-sgmmr    Pulling image "190219044/httpserver:2.0.2"
0s          Normal    Pulled                    pod/cloudnative-httpserver-69fffbc5c7-sgmmr    Successfully pulled image "190219044/httpserver:2.0.2" in 33.290581032s
0s          Normal    Created                   pod/cloudnative-httpserver-69fffbc5c7-sgmmr    Created container httpserver
0s          Normal    Started                   pod/cloudnative-httpserver-69fffbc5c7-sgmmr    Started container httpserver
```

4. 模拟外部访问
```shell
$ curl -ipv4 --resolve foo.bar.com:80:192.168.130.63 http://foo.bar.com/

* Connection #0 to host foo.bar.com left intact
{"hello":"world"}

$ curl http://192.168.130.63:31081/hello

{"hello":"world"}
```

5. 创建 TLS 证书 及给 Ingress 增加 TLS 配置
```
$ openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -out foo-ingress-tls.crt \
    -keyout foo-ingress-tls.key \
    -subj "/CN=foo.bar.com/O=foo-ingress-tls"

$ kubectl create secret tls foo-ingress-tls \
    --key foo-ingress-tls.key \
    --cert foo-ingress-tls.crt

$ kubectl apply -f deployment.yaml

$ curl --insecure -ipv4 --resolve foo.bar.com:443:192.168.130.63 https://foo.bar.com/

< strict-transport-security: max-age=31536000; includeSubDomains
strict-transport-security: max-age=31536000; includeSubDomains

<
* Connection #0 to host foo.bar.com left intact
{"hello":"world"}
```


# 模块二作业

群昵称：我来也-武汉

## 要求
1⃣️  9.25课后作业<br>
内容：编写一个 HTTP 服务器，大家视个人不同情况决定完成到哪个环节，但尽量把1都做完

1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回200

截止时间：10月7日晚23:59前<br>
提示💡：<br>
1. 自行选择做作业的地址，只要提交的链接能让助教老师打开即可
2. 自己所在的助教答疑群是几组，提交作业就选几组

## 说明
偷个懒，直接使用`gin`框架完成作业。

* 在初始化函数`init`中读取环境变量`VERSION`的值。
* 使用中间件，完成`Version`头的注入。
* 注册`/healthz`路由，返回 200 状态码。
* 给所有未注册的路由，指定到函数`HandleGetAllData`中。
* 使用`gin`默认的中间件，打印客户端 IP HTTP 返回码 等信息。

## 使用
1. 本地直接启动
```shell
$ make run

export VERSION=v1.2.3
go run ./
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /healthz                  --> main.main.func3 (5 handlers)
[GIN-debug] Listening and serving HTTP on :8081
```

2. 本地使用`curl`调试
```shell
$ make test

curl -v http://localhost:8081/health
* Uses proxy env variable http_proxy == 'http://127.0.0.1:6152'
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 6152 (#0)
> GET http://localhost:8081/health HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.64.1
> Accept: */*
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 OK
< Accept: */*
< Content-Type: application/json; charset=utf-8
< User-Agent: curl/7.64.1
< Version: 1.2.3
< Date: Sun, 26 Sep 2021 14:51:49 GMT
< Content-Length: 17
<
* Connection #0 to host 127.0.0.1 left intact
{"hello":"world"}* Closing connection 0
curl -v -H "Custom-Header-Key: custom-value" \
		-H 'Current-Header-Array: [ "value1", "value2", "value3" ]' \
		-H 'Current-Header-Map: {"key", "value"}' \
		http://localhost:8081/abc
* Uses proxy env variable http_proxy == 'http://127.0.0.1:6152'
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 6152 (#0)
> GET http://localhost:8081/abc HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/7.64.1
> Accept: */*
> Proxy-Connection: Keep-Alive
> Custom-Header-Key: custom-value
> Current-Header-Array: [ "value1", "value2", "value3" ]
> Current-Header-Map: {"key", "value"}
>
< HTTP/1.1 200 OK
< Accept: */*
< Content-Type: application/json; charset=utf-8
< Current-Header-Array: [ "value1", "value2", "value3" ]
< Current-Header-Map: {"key", "value"}
< Custom-Header-Key: custom-value
< User-Agent: curl/7.64.1
< Version: 1.2.3
< Date: Sun, 26 Sep 2021 14:51:49 GMT
< Content-Length: 17
<
* Connection #0 to host 127.0.0.1 left intact
{"hello":"world"}* Closing connection 0
```

3. 构建镜像并执行
```shell
$ make build run-docker

docker build --build-arg version=v2.1.0 . -t httpserver:latest
...

docker run -it --name=httpserver --rm -p 8081:8081 httpserver:latest
...
```

4. 测试镜像
```shell
$ make test
```

# 模块三作业

群昵称：我来也-武汉

## 要求
- 构建本地镜像。
- 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化（请思考有哪些最佳实践可以引入到 Dockerfile 中来）。
- 将镜像推送至 Docker 官方镜像仓库。
- 通过 Docker 命令本地启动 httpserver。
- 通过 nsenter 进入容器查看 IP 配置。

作业需编写并提交 Dockerfile 及源代码。

## 说明
本次作业大部分在上次就完成了，这次只补充一下`推送镜像`和`nsenter`查看容器IP配置的部分。

### 推送镜像
模拟推送，并不会真的推送成功。😁
```
$ make push

docker push httpserver:latest
```

### nsenter 进入容器查看 IP 配置
```
# 获取容器的进程号
$ PID=$(docker inspect --format {{.State.Pid}} httpserver)

# 进入 nsenter 容器 (MacOS 环境无法直接使用 nsenter 命令，Linux 环境可以跳过该步骤)
$ docker run -it --rm --privileged --pid=host justincormack/nsenter1

# 查看 IP 信息 (需要将 $PID 替换成对应的值)
$ nsenter --target $PID --mount --uts --ipc --net --pid ip addr

1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: tunl0@NONE: <NOARP> mtu 1480 qdisc noop state DOWN qlen 1000
    link/ipip 0.0.0.0 brd 0.0.0.0
3: ip6tnl0@NONE: <NOARP> mtu 1452 qdisc noop state DOWN qlen 1000
    link/tunnel6 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00 brd 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00
10: eth0@if11: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
```

