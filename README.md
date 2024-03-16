# 建设中 🔨

## 开发环境

IDE🧑‍💻： [GoLand](https://www.jetbrains.com/go/)

OS🪟🐧：[Ubuntu 22.04.3 LTS (WSL2)](https://ubuntu.com/desktop/wsl)

```shell
# 环境配置
# Golang
wget https://golang.google.cn/dl/go1.22.1.linux-amd64.tar.gz
sudo tar xfz go1.22.1.linux-amd64.tar.gz -C /usr/local
sudo vim /etc/profile
# export GOROOT=/usr/local/go
# export GOPATH=$HOME/gowork
# export GOBIN=$GOPATH/bin
# export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
source /etc/profile
go version
go env -w GOPROXY="https://goproxy.cn"
go env -w GO111MODULE=on

# Docker
# 
git clone https://github.com/hcjjj/webook.git
```


## 社交平台

**基本介绍**

* 用户登录服务 😺
* 用户关系服务 🧩
* 发帖服务 📃
* 支付服务 💰
* 搜索服务 🔍
* 即时通讯 💬
* Feed 流 🏄

> 如何启动前端：在 webook-fe 目录下先 `npm install` 后 `npm run dev`

**项目结构**

* 参考 [Kratos](https://go-kratos.dev/)、[go-zero](https://go-zero.dev/) 、[Domain-Driven Design](https://zhuanlan.zhihu.com/p/91525839)
* Service - Repository - DAO (Data Access Object) 三层结构 
  * service：领域服务（domain service），一个业务的完整处理过程
  * repository：领域对象的存储，存储数据的抽象
    * dao：数据库操作
  * domain：领域对象
* handler（和HTTP打交道） → service（主要业务逻辑） → repository（数据存储抽象） → dao（数据库操作）

## 技术栈

**第三方库**

* [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web 框架
  * [Middleware](https://github.com/gin-gonic/contrib) - Collection of middlewares created by the community
  * [cors](https://github.com/gin-contrib/cors) -  Official *cross-origin resource sharing* (CORS) gin's middleware
  * [sessions](https://github.com/gin-contrib/sessions) - Gin middleware for session management
  
* [dlclark/regexp2](https://github.com/dlclark/regexp2) - full-featured 正则表达式
* [go-gorm/gorm](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
  * [go-gorm/mysql](https://github.com/go-gorm/mysql) - GORM mysql driver
* [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - Golang implementation of JSON Web Tokens (JWT)

**相关组件**

* [Docker](https://www.docker.com/)
  * [mysql](https://hub.docker.com/_/mysql)
  * [redis](https://hub.docker.com/_/redis)
* [kubernates](https://kubernetes.io/)
  * [kubectl](https://kubernetes.io/docs/tasks/tools/)

## 技术要点
* 用户登录服务
  * 注册、密码加密存储
  * 登录、登录态校验
    * Cookie + Session
    * Session 存储基于 Redis 实现（多实例部署环境）
      * 但是每次请求都要访问 Redis，性能瓶颈问题
      * 换为 JWT（JSON Web Token）机制
        * 这边有个问题需要解决，多实例部署的退出登录功能
    *   刷新登录状态
      * 在登录校验处执行相关逻辑
      * 控制 Session 的有效期
      * 生成一个新的 Token
  * 保护登录系统
    * 限流（限制每个用户每秒最多发送固定数量的请求  ）
      * 基于 Redis 的 IP 限流
    * 增强登录安全
      * 利用 User-Agent 增强安全性  
* 用户关系服务
* 发帖服务
* 支付服务
* 搜索服务 
* 即时通讯
* Feed 流 

**用 Kubernetes 部署 Web 服务器**

交叉编译为 Linux 平台的应用程序 `GOOS=linux GOARCH=amd64 go build -o webook .`

