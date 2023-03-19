# socks5Proxy

[![GitHub license](https://img.shields.io/github/license/guofuwei/socks5proxy)](https://github.com/guofuwei/socks5Proxy/blob/master/LICENSE)[![language](https://img.shields.io/badge/language-Go-blue.svg)](https://golang.org/)![tool](https://img.shields.io/badge/tool-proxy-red.svg)

用golang实现的[socks5协议](./docs/socks5.md)的代理转发，在数据混淆部分使用AES加密算法，在client和server端的通信部分使用自定义的通信协议

## 使用说明

* 在Release中下载对应的server和client文件

* 在Google或者Edge浏览器中安装`SwitchyOmega`扩展插件
* 选择socks5协议，地址和端口填写本机client端的地址和端口
* 在远程服务器上启用server端，在本地启动client端即可开始使用

## 参数说明

### 客户端(client)

```shell
# 默认：客户端监听127.0.0.1:8080，服务器地址127.0.0.1:8081
-local string #设置本地监听地址
   Input listen address(Default 127.0.0.1:8080)
-dest string #设置服务器地址
   Input remote server address(Default 127.0.0.1:8081)
```

### 服务端(server)

```shell
# 默认：服务器监听127.0.0.1:8081
-local string #设置本地监听地址
   Input listen address(Default 127.0.0.1:8081)
```

### 示例

```shell
# 客户端监听在localhost:10000接受浏览器请求，服务器地址设置为45.17.78.90:8080
./client -dest 45.17.78.90:8080 -local localhost:10000
# 服务器监听在所有的8080端口（开放外网连接）
./server -local :8080
```

## 开发说明

### 克隆当前仓库

```shell
git clone https://github.com/guofuwei/socks5Proxy.git
```
### 文件结构

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



