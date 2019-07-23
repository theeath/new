package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"news/models"
	"path"
	"strconv"
	"time"
)

type ArticleController struct {
	beego.Controller
}
////处理下拉框
//func (this*ArticleController)HandleSelect()  {
//	typeName:=this.GetString("select")
//	//处理数据
//	if typeName==""{
//		beego.Info("下拉框传递数据失败")
//		return
//	}
//	o:=orm.NewOrm()
//	var articles[]models.Article
//	o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
//	beego.Info(articles)
//}
func (this*ArticleController)ShowArticleList()  {
	o:=orm.NewOrm()
	qs:=o.QueryTable("Article")
	var articles[] models.Article
	//qs.All(&articles)//select * from Article
	//pageIndex1:=1
	pageIndex:=this.GetString("pageIndex")
	pageIndex1,err:=strconv.Atoi(pageIndex)
	if err!=nil {
		pageIndex1 = 1
	}
	typeName:=this.GetString("select")
	count,err:=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()//返回数据条目数
	pageSize:=2
	//分页
	//pageIndex:=1
	start:=pageSize*(pageIndex1-1)
	qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)//1.pagesize start
	pageCount :=float64(count)/float64(pageSize)
	pageCount1:=math.Ceil(pageCount)
	if err!=nil {
		beego.Info("查询错误")
		return
	}
	FirstPage:=false
	if pageIndex1 == 1 {
		FirstPage=true
	}
	EndPage:=false
	if pageIndex1==int(pageCount1) {
		EndPage = true
	}

	//获取类型数据
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]  = types

	var articlewithtype []models.Article
	//处理数据
	if typeName==""{
		beego.Info("下拉框传递数据失败")
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articlewithtype)//1.pagesize start
	}else {
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articlewithtype)

	}
	userName:=this.GetSession("userName")
	this.Data["userName"] = userName
	this.Data["typeName"] = typeName

	this.Data["FirstPage"] = FirstPage
	this.Data["EndPage"] = EndPage
	this.Data["pageCount"] = pageCount1
	this.Data["count"] = count
	this.Data["pageIndex"] = pageIndex1
	this.Data["articles"] = articlewithtype
	this.Layout = "layout.html"
	this.TplName = "index.html"
}
func (this*ArticleController)ShowAddArticle() {
	//查询类型数据，传递到视图中
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "add.html"

}
func (this*ArticleController)HandleAddArticle()  {
	//那标题
	artiName:=this.GetString("articleName")
	artiContent:=this.GetString("content")
	f,h,err:=this.GetFile("uploadname")
	defer f.Close()
	//判断文件格式
	ext:=path.Ext(h.Filename)
	if ext!=".jpg"&& ext!=".png" && ext!=".jpeg"{
		beego.Info("格式错误")
		return
	}
	//判断文件大小
	if h.Size>5000000 {
		beego.Info("文件太大，不允许上传")
		return
	}
	//不能重名
	fileName:=time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	if err!=nil {
		beego.Info("文件上传失败")
		return
	}
	//插入数据
	//1.获取orm对象
	o:=orm.NewOrm()
	//2.创建一个插入对象
	article:=models.Article{}
	//3.赋值
	article.Content = artiContent
	article.Title = artiName
	article.Img = "/static/img/"+fileName+ext
	//给article对象赋值
	//获取下拉框传递的数据
	typeName:=this.GetString("select")
	if typeName=="" {
		beego.Info("下拉框数据错误")
	}
	//获取type对象
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err=o.Read(&artiType,"TypeName")
	if err!=nil {
		beego.Info("获取类型失败")
		return
	}
	article.ArticleType = &artiType
	//4.插入
	_,err=o.Insert(&article)
	//beego.Info(article.Title,article.Content,article.Id)
	if err!=nil {
		beego.Info("插入数据失败",err)
		return
	}


	this.Redirect("/Article/ShowArticle",302)
}
func (this*ArticleController)ShowArticleDetail()  {
	//获取数据
	id,err:=this.GetInt("articleId")
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
	}
	//操作数据
	o := orm.NewOrm()

	var article models.Article
	article.Id = id

	o.Read(&article)

	//修改阅读量
	article.Count += 1
	//artile:=models.Article{Id:id}
	//多对多插入
	m2m:=o.QueryM2M(&article,"Users")
	userName:=this.GetSession("userName")
	user:=models.User{}
	user.UserName=userName.(string)
	o.Read(&user,"userName")
	_,err=m2m.Add(&user)
	if err!=nil {
		beego.Info("插入失败")
		return
	}
	o.Update(&article)

//多对多查询
	//o.LoadRelated(&article,"Users")
	var users[]models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)







	//返回视图页面
	this.Data["users"] = users
	this.Data["article"] = article
	this.Layout="layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["contentHead"] = "head.html"
	this.TplName = "content.html"

}
//删除
func (this*ArticleController)HandleDelete()  {
	//获取数据
	id,err:=this.GetInt("articleId")
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
	}
	//操作数据
	o := orm.NewOrm()

	var article models.Article
	article.Id = id
	o.Delete(&article)
	this.Redirect("/Article/ShowArticle",302)

}
func (this*ArticleController)ShowUpdate()  {
	//获取数据
	id,err:=this.GetInt("articleId")
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
	}
	//操作数据
	o := orm.NewOrm()

	var article models.Article
	article.Id = id

	err=o.Read(&article)
	if err!=nil {
		beego.Info("查询为空")
		return
	}

	//返回视图页面
	this.Data["article"] = article
	this.TplName = "update.html"

}
func (this*ArticleController)HandleUpdate()  {
	name:=this.GetString("articleName")
	content:=this.GetString("content")
	//获取数据
	id,err:=this.GetInt("articleId")
	//数据校验
	if err != nil{
		beego.Info("传递的链接错误")
	}
	if name==""||content=="" {
		beego.Info("更新数据失败")
		return
	}
	f,h,err:=this.GetFile("uploadname")
	if err!=nil {
		beego.Info("上传文件失败")
		return
	}
	defer f.Close()
	if h.Size >5000000 {
		beego.Info("图片太大")
		return
	}
	ext:=path.Ext(h.Filename)
	if ext!=".jpg"&&ext!=".png"&&ext!=".jpeg" {
		beego.Info("上传文件类型错误")
		return
	}
	filename:=time.Now().Format("2006-01-02 15:04:05")
	this.SaveToFile("uploadname","./static/img/"+filename+ext)
	o:=orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err!=nil {
		beego.Info("要更新文章不存在")
		return
	}
	article.Title = name
	article.Content = content
	article.Img = "./static/img/"+filename+ext
	_,err=o.Update(&article)
	if err!=nil{
		beego.Info("更新失败")
		return
	}
	this.Redirect("/Article/ShowArticle",302)

}
func (this*ArticleController)ShowAddArticleType()  {
	//读取类型吧，
	o:=orm.NewOrm()
	var artiTypes[] models.ArticleType
	_,err:=o.QueryTable("ArticleType").All(&artiTypes)
	if err!=nil {
		beego.Info("查询类型错误")
	}
	this.Data["types"] = artiTypes

	this.TplName = "addType.html"
}
//处理添加类型业务
func (this*ArticleController)HandleAddType()  {
	//获取数据
	typename:=this.GetString("typeName")
	//判断数据
	if typename=="" {
		beego.Info("数据为空")
		return
	}
	//执行插入操作
	o:=orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	_,err:=o.Insert(&artiType)
	if err!=nil {
		beego.Info("插入失败")
		return
	}
	//展示视图
	this.Redirect("/Article/AddArticleType",302)
}
func (this*ArticleController) Logout() {
	//删除登录状态
	this.DelSession("userName")
	//跳转登录页面
	this.Redirect("/",302)
}