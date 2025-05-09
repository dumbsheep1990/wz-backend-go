package main

import (
	"flag"

	"wz-backend-go/internal/admin/config"
	"wz-backend-go/internal/admin/handler"
	"wz-backend-go/internal/admin/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/admin.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting admin server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
