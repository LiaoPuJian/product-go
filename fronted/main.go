package main

import (
	"context"
	"fmt"
	"product-go/common"
	"product-go/fronted/middleware"
	"product-go/fronted/web/controllers"
	"product-go/repositories"
	"product-go/services"
	"time"

	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12/sessions"

	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	tpl := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tpl)

	app.HandleDir("/public", "./fronted/web/public")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {
		panic(fmt.Sprintf("Mysql connect error: %s", err))
	}

	//session
	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 60 * time.Minute,
	})

	//请求上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//注册user
	user := repositories.NewUserRepository("user", db)
	userService := services.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess)
	userPro.Handle(new(controllers.UserController))

	//注册product
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(order)

	productParty := app.Party("/product")
	productPro := mvc.New(productParty)
	//注册中间件
	productParty.Use(middleware.AuthConProduct)

	productPro.Register(productService, orderService, db)
	productPro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
