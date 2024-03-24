# 建设中 🔨


## 项目介绍

**开发环境**

IDE🧑‍💻： [GoLand](https://www.jetbrains.com/go/)

OS🪟🐧：[Ubuntu 22.04.3 LTS (WSL2)](https://ubuntu.com/desktop/wsl)

**开发计划**

- [x] 用户登录服务 👤
  - [x] 注册、登录态校验与刷新
  - [x] 保护登录系统
  - [x] 优化登录性能
  - [ ] 短信验证码登录
  - [ ] 微信扫码登录
- [ ] 发帖服务 📃
- [ ] 用户关系服务 🧩
- [ ] 支付服务 💰
- [ ] 搜索服务 🔍
- [ ] 即时通讯 💬
- [ ] Feed 流 🏄

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
* [tencentcloud-sdk-go](https://github.com/TencentCloud/tencentcloud-sdk-go) - Tencent Cloud API 3.0 SDK for Golang
  * [腾讯云 SMS](https://console.cloud.tencent.com/smsv2) 个人用户无法使用短信服务 API
* ~~[shansuma](https://gitee.com/shansuma/sms-sdk-master) - 闪速码 SMS 的 API 接口~~
* [wire](https://github.com/google/wire) - Compile-time Dependency Injection for Go

**相关环境**

* [Node.js](https://nodejs.org/en)
  * 启动前端：在 webook-fe 目录下先 `npm install` 后 `npm run dev`
* [Docker](https://www.docker.com/)
  * [镜像源](https://yeasy.gitbook.io/docker_practice/install/mirror)（还是挂代理方便）
  * [mysql](https://hub.docker.com/_/mysql) - An open-source relational database management system (RDBMS)
  * [redis](https://hub.docker.com/_/redis) - An open-source in-memory storage
* [kubernates](https://kubernetes.io/)
  * [Kubernetes cluster architecture](https://kubernetes.io/docs/concepts/architecture/)
  * [kubectl](https://kubernetes.io/docs/tasks/tools/) - The Kubernetes command-line tool
  * [HELM](https://helm.sh/) - The package manager for Kubernetes
  * [ingress-nignx](https://github.com/kubernetes/ingress-nginx) - Ingress-NGINX Controller for Kubernetes
* [wrk](https://github.com/wg/wrk) - Modern HTTP benchmarking tool
* [ekit](https://github.com/ecodeclub/ekit) - 支持泛型的工具库

## 技术要点
**业务功能**

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
  * 优化登录性能
  * 短信验证码登录
    * 验证码是一个独立的功能 （登录、修改密码、危险操作的二次验证）
    * 短信服务也是独立的（方便更换供应商）
    * 验证码登录功能 → 验证码功能 → 短信服务（最基础的服务）
* 用户关系服务
* 发帖服务
* 支付服务
* 搜索服务 
* 即时通讯
* Feed 流

**编程思想**

* 控制反转（Inversion of Control, IoC）
  * 依赖注入（Dependency Injection）
  * 依赖查找、依赖发现（Go 里面没有）
* 面向接口编程

# 部署应用

**环境配置**

```shell
# Ubuntu 22.04.3 LTS
# Golang
wget https://golang.google.cn/dl/go1.22.1.linux-amd64.tar.gz
sudo tar xfz go1.22.1.linux-amd64.tar.gz -C /usr/local
sudo vim /etc/profile
# export GOROOT=/usr/local/go
# export GOPATH=$HOME/go
# export GOBIN=$GOPATH/bin
# export PATH=$GOPATH:$GOBIN:$GOROOT/bin:$PATH
source /etc/profile
go version
go env -w GOPROXY="https://goproxy.cn"
go env -w GO111MODULE=on

# Docker
# 

# Kubernetes
#
```

**用 Kubernetes 部署 Web 服务**

交叉编译

```shell
# Windows → Linux
# powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o .\build\webook
go build -tags=k8s -o .\build\webook
# Mac → Linux
GOOS=linux GOARCH=amd64 go build -o /build/webook
```

编写 `Dockerfile`

```dockerfile
# 基础镜像
FROM ubuntu:20.04
# 把编译后的打包进这个镜像，放到工作目录 /app
COPY /build/webook /app/webook
WORKDIR /app
# CMD 是执行命令
# 最佳
ENTRYPOINT ["/app/webook"]
```

```shell
# 构建
docker build -t hcjjj/webook:v0.0.1 .
# 删除
docker rmi -f hcjjj/webook:v0.0.1
# 可以将上述命令都写在 Makefile 里面
```

编写 `k8s.yaml` 后

```shell
# 启动 deployment
kubectl apply -f k8s-webook-deployment.yaml
# 查看 
kubectl get deployments
kubectl get pods
# 查看 POD 的日志
kebectl get logs -f webook-5b4c5b9-4g74z
# 启动 services
kubectl apply -f k8s-webook-service.yaml
# 查看
kubectl get services
# 停止
kubectl delete service webook
kubectl delete deployment webook
```

**用 Kubernetes 部署 Mysql**

```shell
# Mysql 持久化
# 启动
kubectl apply -f k8s-mysql-deployment.yaml
kubectl apply -f k8s-mysql-service.yaml
kubectl apply -f k8s-mysql-pv.yaml
kubectl apply -f k8s-mysql-pvc.yaml
# 查看
kubectl get pv
kubectl get pvc
# 停止
kubectl delete service webook-mysql
kubectl delete deployment webook-mysql
kubectl delete pvc webook-mysql-claim
kubectl delete pv webook-mysql-pv
```

**用 Kubernetes 部署 Redis**

```shell
kubectl apply -f k8s-redis-deployment.yaml
kubectl apply -f k8s-redis-service.yaml
kubectl delete service webook-redis
kubectl delete deployment webook-redis
```

**用 Kubernetes 部署 nginx**

```shell
# 本地环境需要修改 host 到 ip 的映射，host 在 k8s-ingress-nginx.yaml 里面
# ❯ ping  hcjjj.webook.com
# PING hcjjj.webook.com (127.0.0.1) 56(84) bytes of data.
# 64 bytes from localhost (127.0.0.1): icmp_seq=1 ttl=64 time=0.028 ms
# 使用 clash for windows 的话，同时需要在 Bypass Domain/IPNet 中添加 

# 安装 ingress-nignx 
helm upgrade --install ingress-nginx ingress-nginx  --repo https://kubernetes.github.io/ingress-nginx  --namespace ingress-nginx --create-namespace
# 查看
kubectl get service --namespace ingress-nginx
# 启动
kubectl apply -f k8s-ingress-nginx.yaml
# 停止
kubectl get ingresses
kubectl delete ingress webook-ingress
kubectl delete namespaces ingress-nginx
```
