package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

type AttachController struct {
	beego.Controller
}

func (c *AttachController) Get() {
	// RequestURI 中如果有中文会进行编码，要想获得文件名需要反编码
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	defer f.Close()
	_, err = io.Copy(c.Ctx.ResponseWriter, f)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
}
