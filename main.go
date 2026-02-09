package main

import (
	"flag"
	"fmt"

	"accesscontrol/internal/config"
	"accesscontrol/internal/errorx"
	"accesscontrol/internal/handler"
	"accesscontrol/internal/svc"
	"net/http"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)
var configFile = flag.String("f", "etc/config.yaml", "the config file")
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 自定义错误处理
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.CodeError:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})
	ctx := svc.NewServiceContext(c)
	handler.SetupRoutes(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
