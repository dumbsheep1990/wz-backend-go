package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"wz-backend-go/internal/config"
	"wz-backend-go/internal/handler/admin"
	"wz-backend-go/internal/middleware"
	"wz-backend-go/internal/svc"
)

var configFile = flag.String("f", "configs/admin.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.AdminConfig
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.Log)
	defer logx.Close()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewAdminServiceContext(c)
	admin.RegisterHandlers(server, ctx)

	// 注册中间件
	server.Use(middleware.AdminCheck)

	fmt.Printf("Starting admin server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
