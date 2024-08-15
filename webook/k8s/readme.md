## 单体版本应用部署

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
# Kubernetes
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

编写 `k8s.yaml`

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

**用 Kubernetes 部署 Nginx**

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

## 