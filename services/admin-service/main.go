package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"wz-backend-go/services/admin-service/config"
	"wz-backend-go/services/admin-service/internal/service"
)

var configFile = flag.String("f", "services/admin-service/config/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)
	defer logx.Close()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := service.NewServiceContext(c)
	server.AddRoutes(server.Routes())

	// 注册中间件
	server.Use(ctx.AdminAuthMiddleware)

	fmt.Printf("Starting admin-service at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
