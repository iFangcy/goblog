package main

import (
	"myblog/models"
	_ "myblog/routers"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	// 第二个参数，默认连接上数据库时不会创建表，这里需要创建表，false 表示当有表示不要删掉重建
	// 第三个参数，打印信息
	orm.RunSyncdb("default", false, true)
	// 创建附件目录
	os.Mkdir("attachment", os.ModePerm)

	// 国际化支持
	i18n.SetMessage("en-US", "conf/locale_en-US.ini")
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")
	// 注册模板函数，让这个函数可以在模板中使用
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.Run()
}
