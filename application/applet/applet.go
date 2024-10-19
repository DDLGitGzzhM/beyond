package main

import (
	"flag"
	"fmt"

	"beyond/pkg/xcode"

	"github.com/zeromicro/go-zero/rest/httpx"

	"beyond/application/applet/internal/config"
	"beyond/application/applet/internal/handler"
	"beyond/application/applet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

// branch
// checkout
// pull
// push
// ---
// git add ,  git commit , git push
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
