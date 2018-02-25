package models

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME        = "data/myblog.db"
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int64
}

// 评论
type Comment struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// 必须有个数据库叫 default
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
}

// 添加分类操作
func AddCategory(name string) error {
	// 获取 orm 对象
	o := orm.NewOrm()
	// 创建一个分类对象
	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}
	// 获取 category 表的 query
	qs := o.QueryTable("category")
	// 查找当前表中是否已经使用了这个分类名称
	err := qs.Filter("title", name).One(cate)
	// 表中查找到数据时 err 为 nil， 没有查找的数据时 err 有值
	if err == nil {
		// err 为 nil 表示表中已经有了这个分类
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		// 有错，插入失败
		return err
	}
	return nil
}

// 获取所有分类
func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// 删除一个分类
func DeleteCategory(id string) error {
	// 将删除按钮获取的字符串 id 转成 10 进制 64 位数值
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	// 删除、读取操作必须指定主键
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

// 添加一篇文章
func AddTopic(title, content, label, category, attachment string) error {
	// 传进来的是空格区分的多个标签
	// 处理后效果： $beego#$bee# 是为了关键字搜索是可以区分相似的字符串
	if len(label) > 0 {
		label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	}
	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Content:    content,
		Category:   category,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
		Labels:     label,
		Attachment: attachment,
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	// 添加成功，修改相应分类文章数量
	if len(category) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", category).One(cate)
		if err == nil {
			cate.TopicCount++
			_, err = o.Update(cate)
		}
	}
	return err
}

// 获取所有文章列表
func GetAllTopics(category, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		// 倒序排列，只在首页使用，需要判断是否根据分类显示
		if len(category) > 0 {
			qs = qs.Filter("category", category)
		}
		if len(label) > 0 {
			// 根据标签获取
			qs = qs.Filter("labels__contains", "$"+label+"#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		// 顺序排列
		_, err = qs.All(&topics)
	}

	if err != nil {
		beego.Error(err)
	}
	return topics, err
}

// 获取指定一篇文章
func GetTopic(tid string) (*Topic, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)
	if err != nil {
		return nil, err
	}
	// 增加浏览数
	topic.Views++
	_, err = o.Update(topic)
	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)
	return topic, nil
}

// 修改一篇旧文章
func ModifyTopic(tid, title, content, category, label, attachment string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		beego.Error(err)
		return err
	}
	if len(label) > 0 {
		label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
	}
	var oldCate, oldAttach string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldAttach = topic.Attachment
		oldCate = topic.Category
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Updated = time.Now()
		topic.Labels = label
		topic.Attachment = attachment
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	// 更新相应分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			if cate.TopicCount > 0 {
				cate.TopicCount--
				_, err = o.Update(cate)
			}
		}
	}

	if len(category) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", category).One(cate)
		if err == nil {
			cate.TopicCount++
			_, err = o.Update(cate)
		}
	}

	// 删除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment", oldAttach))
	}

	return nil
}

// 删除一篇文章
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		beego.Error(err)
		return err
	}

	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		// 数据库中有这篇文章
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}

	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			if cate.TopicCount > 0 {
				cate.TopicCount--
				_, err = o.Update(cate)
			}
		}
	}

	return err
}

// 添加一条评论
func AddReply(tid, nickname, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		beego.Error(err)
		return err
	}
	reply := &Comment{
		Tid:     tidNum,
		Name:    nickname,
		Content: content,
		Created: time.Now(),
	}
	o := orm.NewOrm()
	_, err = o.Insert(reply)
	if err != nil {
		return err
	}
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
	}
	return err
}

// 删除一条评论
func DeleteReply(rid string) error {
	ridNum, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		beego.Error(err)
		return err
	}
	o := orm.NewOrm()
	var tidNum int64
	comment := &Comment{Id: ridNum}
	if o.Read(comment) == nil {
		tidNum = comment.Id
		_, err = o.Delete(comment)
		if err != nil {
			return err
		}
	}
	// 获取这篇文章评论列表，按时间最新在最前排列
	replies := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}

	// 更新文章评论数，最新评论时间
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		if len(replies) > 0 {
			topic.ReplyTime = replies[0].Created
			topic.ReplyCount = int64(len(replies))
			_, err = o.Update(topic)
		}
	}

	return err
}

// 获取指定文章所有评论
func GetAllReplies(tid string) ([]*Comment, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	replies := make([]*Comment, 0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidNum).All(&replies)
	return replies, err
}
