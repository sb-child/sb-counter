package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"sb-counter/internal/controller/counter"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// s.Group("/", func(group *ghttp.RouterGroup) {
			// 	group.Middleware(ghttp.MiddlewareHandlerResponse)
			// 	group.Bind(
			// 		hello.New(),
			// 	)
			// })
			s.Group(g.Config().MustGet(ctx, "sbcounter.rootDir").String(), func(group *ghttp.RouterGroup) {
				group.GET("/:user_path/:mode/:output", counter.Handler)
			})
			s.Run()
			return nil
		},
	}
)
