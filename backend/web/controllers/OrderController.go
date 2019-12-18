package controllers

import (
	"fmt"
	"product-go/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Crx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orders, err := o.OrderService.GetAllOrderInfo()
	if err != nil {
		fmt.Println(err)
		o.Crx.Application().Logger().Debug("查询订单失败")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"products": orders,
		},
	}
}
