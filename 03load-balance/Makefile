COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))

# go output dir
GO_OUT_DIR := $(ROOT_DIR)/go/rpc

# Protobuf 文件存放路径
APIROOT=$(ROOT_DIR)/proto

.PHONY: echo
echo:
	@echo $(APIROOT)
	@echo $(ROOT_DIR)

.PHONY: go-protoc
go-protoc:
	@protoc -I./proto --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
	hello.proto

.PHONY: node-protoc
node-protoc:
	@protoc -I./proto --js_out=import_style=commonjs,binary:./node \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./node hello.proto