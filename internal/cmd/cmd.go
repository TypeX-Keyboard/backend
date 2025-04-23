package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"keyboard-api-go/internal/controller/bot"
	"keyboard-api-go/internal/controller/chain"
	"keyboard-api-go/internal/controller/secure"
	"keyboard-api-go/internal/controller/system"
	"keyboard-api-go/internal/controller/token"
	"keyboard-api-go/internal/controller/wallet"
	"keyboard-api-go/internal/packed/ws"
	"keyboard-api-go/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			MustInit(ctx)
			//启动服务
			ws.StartWebSocket(ctx)
			s := g.Server()
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				//注册路由
				group.ALL("/socket", ws.WsHandle)
				group.Middleware(service.Middleware().CORS)
				group.Middleware(service.Middleware().Context)
				group.Middleware(service.Middleware().Decrypted)
				group.Middleware(service.Middleware().Encrypted)
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 安全模块
				group.Group("/secure", func(group *ghttp.RouterGroup) {
					group.Bind(
						secure.NewV1(),
					)
				})
				// 区块链模块
				group.Group("/chain", func(group *ghttp.RouterGroup) {
					group.Bind(
						chain.NewV1(),
					)
				})
				// 钱包模块
				group.Group("/wallet", func(group *ghttp.RouterGroup) {
					group.Bind(
						wallet.NewV1(),
					)
				})
				// token模块
				group.Group("/token", func(group *ghttp.RouterGroup) {
					group.Bind(
						token.NewV1(),
					)
				})
				// bot模块
				group.Group("/bot", func(group *ghttp.RouterGroup) {
					group.Bind(
						bot.NewV1(),
					)
				})
				// system模块
				group.Group("/system", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewV1(),
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
