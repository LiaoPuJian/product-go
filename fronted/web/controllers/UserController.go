package controllers

import (
	"product-go/common"
	"product-go/encrypt"
	"product-go/models"
	"product-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (u *UserController) PostRegister() {
	var (
		nickName = u.Ctx.FormValue("nickName")
		userName = u.Ctx.FormValue("userName")
		password = u.Ctx.FormValue("password")
	)

	user := &models.User{
		NickName: nickName,
		UserName: userName,
		Password: password,
	}

	_, err := u.Service.AddUser(user)
	u.Ctx.Application().Logger().Debug(err)

	if err != nil {
		u.Ctx.Redirect("/user/error")
	}
	u.Ctx.Redirect("/user/login")
	return
}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	var (
		userName  = u.Ctx.FormValue("userName")
		passsword = u.Ctx.FormValue("password")
	)

	//验证账号密码是否正确
	user, err := u.Service.IsPwdSuccess(userName, passsword)
	u.Ctx.Application().Logger().Debug(err)
	if err != nil {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	//将用户id写入到cookie中
	uidStr := strconv.FormatInt(user.ID, 10)
	common.GlobalCookie(u.Ctx, "uid", uidStr)

	//将用户id加密，并将加密后的sign写入到Cookie中
	sign, err := encrypt.EnPwdCode([]byte(uidStr))
	if err != nil {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	common.GlobalCookie(u.Ctx, "sign", sign)

	return mvc.Response{
		Path: "/product/detail",
	}
}
