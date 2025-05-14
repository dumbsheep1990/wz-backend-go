@echo off
echo 启动万知市场站点建设器后端服务...

echo 构建服务...
cd /d %~dp0\..
go build -o bin\site-service.exe services\site-service\main.go
go build -o bin\page-service.exe services\page-service\main.go
go build -o bin\component-service.exe services\component-service\main.go
go build -o bin\render-service.exe services\render-service\main.go
go build -o bin\api-gateway.exe api-gateway\main.go

echo 启动服务...
start cmd /k "cd /d %~dp0\.. && bin\site-service.exe"
timeout /t 2
start cmd /k "cd /d %~dp0\.. && bin\page-service.exe"
timeout /t 2
start cmd /k "cd /d %~dp0\.. && bin\component-service.exe"
timeout /t 2
start cmd /k "cd /d %~dp0\.. && bin\render-service.exe"
timeout /t 2
start cmd /k "cd /d %~dp0\.. && bin\api-gateway.exe"

echo 所有服务已启动
echo API网关地址: http://localhost:8080 