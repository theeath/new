package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"news/models"
	"time"
)

type RegController struct {
	beego.Controller
}

func (this*RegController)ShowReg()  {
	this.TplName = "register.html"
}
/*
	状态码
1xx 继续发送
2xx 请求成功 200
3xx 资源转移 302 重定向
4xx 请求错误 404
5xx 服务器错误
c.tplname
1.服务器端的功能，浏览器没有再次发送请求，用上次的地址，地址不变
2.传递数据
redirect
1.浏览器端的请求，浏览器再次访问了服务器，地址改变
2.不能传递数据
	1.拿到浏览器传递的数据
	2.数据处理
	3.插入数据库
	4.返回视图
*/
func (this*RegController)HandleReg()  {
	//1.拿到浏览器传递的数据
	name:=this.GetString("userName")
	passwd:=this.GetString("password")
	//2.数据处理
	if name==""||passwd == "" {
		beego.Info("用户名密码不能为空")
		this.TplName = "register.html"
		return
	}
	//插入数据
	o:=orm.NewOrm()
	user := models.User{}
	user.UserName = name
	user.Passwd = passwd
	_,err:=o.Insert(&user)
	if err!=nil {
		beego.Info("插入数据失败")
		return

	}
	this.Redirect("/",302)
	//this.Ctx.WriteString("注册成功")
}

type LoginController struct {
	beego.Controller
}

func (this*LoginController)ShowLogin()  {
	name:=this.Ctx.GetCookie("userName")
	if name!=""{
		this.Data["userName"] = name
		this.Data["check"] = "checked"
	}
	this.TplName = "login.html"
}
func (this*LoginController)HandleLogin()  {
	name:=this.GetString("userName")
	passwd:=this.GetString("password")
	//beego.Info(name,passwd)
	if name==""||passwd == "" {
		beego.Info("用户名密码不能为空")
		this.TplName = "login.html"
		return
	}
	//查找数据
	o:=orm.NewOrm()
	user:=models.User{}
	user.UserName = name
	err:=o.Read(&user,"UserName")
	if err!=nil {
		beego.Info("用户名失败")
		this.TplName = "login.html"
		return
	}
	//判断密码
	if user.Passwd!=passwd {
		beego.Info("密码错误")
		this.TplName = "login.html"
		return
	}
	check:=this.GetString("remember")
	if check == "on" {
		this.Ctx.SetCookie("userName",name,time.Second*3600)

	}else {
		this.Ctx.SetCookie("userName","sss",-1)
	}
	this.SetSession("userName",name)
	//this.Ctx.WriteString("登录成功")
	this.Redirect("/Article/ShowArticle",302)

}