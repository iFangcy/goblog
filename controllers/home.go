package controllers

import (
	"myblog/models"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
}

func (c *BaseController) Prepare() {
	lang := c.GetString("lang")
	if lang == "zh-CN" {
		c.Lang = lang
	} else {
		c.Lang = "en-US"
	}
	c.Data["Lang"] = c.Lang
}

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.TplName = "home.html"
	c.Data["IsHome"] = true

	// 根据 cookie 中的用户名和密码判断当前用户是否已经登录
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	cate := c.Input().Get("cate")
	label := c.Input().Get("label")
	topics, err := models.GetAllTopics(cate, label, true)
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"] = topics

	categories, err := models.GetAllCategories()
	if err != nil {
		beego.Error(err)
		return
	}
	c.Data["Categories"] = categories
	// 国际化，可以这里处理，也可以在 main 中注册模板函数后在模板中处理
	c.Data["Hi"] = c.Tr("hi")
	c.Data["Bye"] = c.Tr("bye")
	// 这里通过模板中处理国际化
	// c.Data["Hi"] = "hi"
	// c.Data["Bye"] = "bye"
}
