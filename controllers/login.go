package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	isExit := c.Input().Get("exit") == "true"
	if isExit {
		c.Ctx.SetCookie("username", "", -1, "/")
		c.Ctx.SetCookie("password", "", -1, "/")
		c.Redirect("/", 302)
		return
	}
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	username := c.Input().Get("username")
	password := c.Input().Get("password")
	autoLogin := c.Input().Get("autoLogin") == "on"

	if username == beego.AppConfig.String("username") && password == beego.AppConfig.String("password") {
		// 设置 cookie 有效时间，默认关闭浏览器就删除
		maxAge := 0
		if autoLogin {
			// 如果设置了自动登录，可以设置一个比较大的值
			maxAge = 1<<31 - 1
		}
		// 设置 cookie
		c.Ctx.SetCookie("username", username, maxAge, "/")
		c.Ctx.SetCookie("password", password, maxAge, "/")
	}
	c.Redirect("/", 302)

	return
}

func checkAccount(ctx *context.Context) bool {
	// 获取 cookie 和数据库中用户名、密码比对
	// 这有个坑爹的地方，这里的 context 使用到的包 github.com/astaxie/beego/context 不会自己导入
	ck, err := ctx.Request.Cookie("username")
	if err != nil {
		return false
	}
	username := ck.Value

	ck, err = ctx.Request.Cookie("password")
	if err != nil {
		return false
	}
	password := ck.Value
	return username == beego.AppConfig.String("username") && password == beego.AppConfig.String("password")
}
