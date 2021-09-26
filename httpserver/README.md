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

docker run -it --rm -p 8081:8081 httpserver:latest
...
```

4. 测试镜像
```shell
$ make test
```

