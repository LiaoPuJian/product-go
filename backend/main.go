package main

import (
	"context"
	"product-go/backend/web/controllers"
	"product-go/common"
	"product-go/repositories"
	"product-go/services"

	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12"
)

func main() {
	//1.创建iris实例
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello!</b>")
	})

	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板目录和布局文件
	tmp := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").Reload(true)
	//注册
	app.RegisterView(tmp)
	//4.设置模板目录
	app.HandleDir("/assets", "./backend/web/assets")
	//出现问题跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//连接mysql
	db, err := common.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	//上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//5.注册控制器
	//数据库映射
	productRepository := repositories.NewProductManager("product", db)
	//service
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}

func RegisterController(name string) {

}
