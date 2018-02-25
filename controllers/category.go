package controllers

import (
	"myblog/models"

	"github.com/astaxie/beego"
)

type CategoryController struct {
	beego.Controller
}

func (c *CategoryController) Get() {
	c.Data["IsCategory"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	// 获取所有分类
	cates, err := models.GetAllCategories()
	if err == nil {
		c.Data["Categories"] = cates
	} else {
		beego.Error(err)
	}
	c.TplName = "category.html"

	// 获取分类页面的操作，添加分类，删除分类
	op := c.Input().Get("op")
	switch op {
	case "add":
		// 添加分类操作
		name := c.Input().Get("name")
		if len(name) == 0 {
			break
		}
		err := models.AddCategory(name)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return
	case "del":
		// 删除分类操作
		id := c.Input().Get("id")
		if len(id) == 0 {
			break
		}
		err := models.DeleteCategory(id)
		if err != nil {
			beego.Error(err)
		}
		c.Redirect("/category", 302)
		return
	}

}
