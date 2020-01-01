package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"product-go/models"
	"product-go/services"
	"strconv"

	"github.com/kataras/iris/v12/mvc"

	"github.com/kataras/iris/v12"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	DB             *sql.DB
}

//获取产品详情页
func (p *ProductController) GetDetail() mvc.View {
	//这里先写死默认取id为1的数据
	product, err := p.ProductService.GetProductById(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	productInfo := p.Ctx.URLParam("productID")
	userInfo := p.Ctx.GetCookie("uid")
	pid, err := strconv.Atoi(productInfo)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductById(int64(pid))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	var orderId int64
	showMessage := "抢购失败！"

	if product.ProductNum > 0 {

		//这里应该使用事务  但是如下写法事务不会生效，因为事务应该使用tx，不应该走另外的db
		tx, _ := p.DB.Begin()
		product.ProductNum -= 1
		errUpdateProduct := p.ProductService.UpdateProduct(product)

		//创建订单
		userId, _ := strconv.Atoi(userInfo)

		order := &models.Order{
			UserId:      int64(userId),
			ProductId:   product.ID,
			OrderStatus: models.OrderSuccess,
		}

		var errOrderInsert error
		orderId, errOrderInsert = p.OrderService.InsertOrder(order)
		errOrderInsert = errors.New("插入报错了！")
		if errUpdateProduct != nil || errOrderInsert != nil {
			tx.Rollback()
			fmt.Println("出现错误，事务回滚")
			p.Ctx.Application().Logger().Debug(err)
		} else {
			tx.Commit()
			showMessage = "抢购成功！"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderId,
			"showMessage": showMessage,
		},
	}

}
