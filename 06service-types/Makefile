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
	@echo $(shell find $(APIROOT) -name "*.proto")

.PHONY: go-protoc
go-protoc:
	@protoc -I$(APIROOT) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
	$(shell find $(APIROOT) -name "*.proto")
