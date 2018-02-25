package controllers

import (
	"myblog/models"

	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

// 添加一条评论
func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	nickname := c.Input().Get("nickname")
	content := c.Input().Get("content")

	err := models.AddReply(tid, nickname, content)
	if err != nil {
		beego.Error(err)
		return
	}
	c.Redirect("/topic/view/"+tid, 302)
}

// 删除一条评论
func (c *ReplyController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	rid := c.Input().Get("rid")
	tid := c.Input().Get("tid")
	err := models.DeleteReply(rid)
	if err != nil {
		beego.Error(err)
		return
	}
	c.Redirect("/topic/view/"+tid, 302)
}
