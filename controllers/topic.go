package controllers

import (
	"myblog/models"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	topics, err := models.GetAllTopics("", "", false)
	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Topics"] = topics
	}
	c.TplName = "topic.html"
}

// 跳到提交文章页面
func (c *TopicController) Add() {
	c.TplName = "topic_add.html"
}

// 提交一篇文章
func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	tid := c.Input().Get("tid")
	category := c.Input().Get("category")
	label := c.Input().Get("label")

	// 判断是否上传了附件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		// 有附件
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(tid) == 0 {
		// 添加新文章
		err = models.AddTopic(title, content, label, category, attachment)
	} else {
		// 修改旧文章
		err = models.ModifyTopic(tid, title, content, category, label, attachment)
	}
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
}

// 浏览一篇文章
func (c *TopicController) View() {
	c.TplName = "topic_view.html"
	topic, err := models.GetTopic(c.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	c.Data["Topic"] = topic
	c.Data["Labels"] = strings.Split(topic.Labels, " ")
	c.Data["Tid"] = c.Ctx.Input.Param("0")
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	replies, err := models.GetAllReplies(c.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		return
	}
	c.Data["Replies"] = replies
}

// 修改一篇文章
func (c *TopicController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "topic_modify.html"
	tid := c.Input().Get("tid")
	topic, err := models.GetTopic(tid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Topic"] = topic
	c.Data["Tid"] = tid
}

// 删除一篇文章
func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	err := models.DeleteTopic(c.Input().Get("tid"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
}
