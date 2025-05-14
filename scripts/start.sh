#!/bin/bash
echo "启动万知市场站点建设器后端服务..."

echo "构建服务..."
cd "$(dirname "$0")/.."
mkdir -p bin
go build -o bin/site-service services/site-service/main.go
go build -o bin/page-service services/page-service/main.go
go build -o bin/component-service services/component-service/main.go
go build -o bin/render-service services/render-service/main.go
go build -o bin/api-gateway api-gateway/main.go

echo "启动服务..."
bin/site-service > logs/site-service.log 2>&1 &
sleep 2
bin/page-service > logs/page-service.log 2>&1 &
sleep 2
bin/component-service > logs/component-service.log 2>&1 &
sleep 2
bin/render-service > logs/render-service.log 2>&1 &
sleep 2
bin/api-gateway > logs/api-gateway.log 2>&1 &

echo "所有服务已启动"
echo "API网关地址: http://localhost:8080" 