package main

import (
	"github.com/orangebees/go-layout/internal/application"
	"github.com/orangebees/go-layout/internal/service"
	"github.com/savsgio/atreugo/v11"
	_ "go.uber.org/automaxprocs"
	"strconv"
)

func main() {

	server := application.GetAtreugo()
	//设置健康检查api GET /health -> ok
	application.SetHealthResponse(server)

	//设置路由组 /api
	api := server.NewGroupPath("/api")

	//设置路由组 /api/v1
	v1 := api.NewGroupPath("/v1")

	//为路由组 /api/v1 设置响应为json类型
	application.SetJsonContentType(v1)

	//为路由组请求与响应设置RequestId
	application.SetRequestId(v1)

	//自动装载会话存储
	application.AutoLoadSaveSessionStore(v1)

	//设置访问日志
	application.SetAccessLog(v1)

	//依赖注入服务 Mock
	application.NewService(service.NewMock())

	//依赖注入服务 Default 使用mysql和redis
	//application.NewService(service.NewDefault(
	//	application.GetSqlxMysqlClient(),
	//	application.GetSqlxMysqlTablePrefix(),
	//	application.GetRedis(),
	//	application.GetRedisKeyPrefix()))
	defer application.Close()
	appService := application.GetService()

	//注册 /api/v1/user GET 请求 控制器，在控制器中处理参数并传入service
	v1.GET("/user", func(ctx *atreugo.RequestCtx) error {
		//来自Query的参数
		idbytes := ctx.QueryArgs().Peek("id")
		id, err := strconv.Atoi(string(idbytes))
		if err != nil {
			application.GetLogger().Err(err).Send()
			return nil
		}
		ctx.SetBodyString(appService.GetById(id))
		return nil
	})
	//注册 /api/v1/user Post 请求 控制器，在控制器中处理参数并传入service
	v1.POST("/user2", func(ctx *atreugo.RequestCtx) error {
		//来自PostArgs
		id, err := ctx.PostArgs().GetUint("id")
		if err != nil {
			application.GetLogger().Err(err).Send()
			return nil
		}
		g := appService.Get()
		if g != nil {

		}
		ctx.SetBodyString(appService.GetById(id))
		return nil
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
