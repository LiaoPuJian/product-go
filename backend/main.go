package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	//1.创建iris实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板目录和布局文件
	tmp := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").Reload(true)
	//注册
	app.RegisterView(tmp)
	//4.设置模板目录
	//出现问题跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//5.注册控制器

	//6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
