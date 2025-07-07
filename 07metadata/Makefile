.PHONY: proto clean run-server run-client

# 生成protobuf代码
proto:
	@echo "生成protobuf Go代码..."
	@mkdir -p rpc
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user.proto
	@echo "protobuf代码生成完成"

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	@rm -rf rpc/
	@echo "清理完成"

# 运行服务端
run-server:
	@echo "启动gRPC服务端..."
	@go run server/main.go

# 运行客户端
run-client:
	@echo "启动gRPC客户端..."
	@go run client/main.go

# 安装依赖
install:
	@echo "安装依赖..."
	@go mod tidy
	@echo "依赖安装完成"
