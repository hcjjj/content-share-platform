.PHONY: docker

docker:
	# 删除旧的
	@rm /build/webook || true
	# 编译前 tidy 复制 go.sum 不对
	@go mod tidy
	# 交叉编译
	@GOOS=linux GOARCH=amd64 go build -o /build/webook
	# 打包为镜像
	@docker rmi -f hcjjj/webook:v0.0.1
	@docker build -t hcjjj/webook:v0.0.1 .