# socks5Proxy

[![GitHub license](https://img.shields.io/github/license/shikanon/socks5proxy)](https://github.com/shikanon/socks5proxy/blob/master/LICENSE)

[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)

用golang实现的socks5协议的代理转发，在数据混淆部分使用AES加密算法，在client和server端的通信部分使用自定义的通信协议

文件结构

```shell
core # socks5Proxy的核心部分
core/aes.go # 引用aes加密算法
core/encryption.go # 实现加密读写和双向加密copy
core/sockauth.go # socks5头的处理
core/sockdest.go # 解析socks5的dest地址
core/sockforward.go # 实现双向加密转发

cmd # 执行文件
cmd/server/main.go # server端启动文件
cmd/client/main.go # client端启动文件

server/server.go # server端主要文件
client/client.go # client端主要文件

config.go # 配置server和client通信的USERNAME和PASSWORD及AES的加密向量
```

* [socks5协议介绍](./docs/socks5.md)

## 使用说明

#### 客户端(client)启动

 在项目根文件下`cd`到`./cmd/client/`运行`go run main.go`，下面为命令行参数

```shell
  -local string #设置本地监听地址
    	Input listen address(Default 127.0.0.1:8080)
  -dest string #设置服务器地址
    	Input remote server address(Default 127.0.0.1:8081)
```

#### 服务端(server)启动

 在项目根文件下`cd`到`./cmd/server/`运行`go run main.go`，下面为命令行参数

```shell
  -local string #设置本地监听地址
    	Input listen address(Default 127.0.0.1:8081)
```



