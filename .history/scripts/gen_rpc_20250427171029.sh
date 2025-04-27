#!/bin/bash

# 生成interaction服务RPC代码
goctl rpc protoc api/rpc/interaction.proto --go_out=./rpc/interaction --go-grpc_out=./rpc/interaction --zrpc_out=./rpc/interaction

# 生成ai服务RPC代码
goctl rpc protoc api/rpc/ai.proto --go_out=./rpc/ai --go-grpc_out=./rpc/ai --zrpc_out=./rpc/ai

# 生成notification服务RPC代码
goctl rpc protoc api/rpc/notification.proto --go_out=./rpc/notification --go-grpc_out=./rpc/notification --zrpc_out=./rpc/notification

# 生成file服务RPC代码
goctl rpc protoc api/rpc/file.proto --go_out=./rpc/file --go-grpc_out=./rpc/file --zrpc_out=./rpc/file

# 生成statistics服务RPC代码
goctl rpc protoc api/rpc/statistics.proto --go_out=./rpc/statistics --go-grpc_out=./rpc/statistics --zrpc_out=./rpc/statistics 