package routers

import (
	"github.com/astaxie/beego/context"
	"news/controllers"
	"github.com/astaxie/beego"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
    beego.InsertFilter("/Article/*",beego.BeforeRouter,FiltFunc)
	beego.Router("/register", &controllers.RegController{},"get:ShowReg;post:HandleReg")
	beego.Router("/", &controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle",&controllers.ArticleController{},"get:ShowArticleList")
	beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	//显示文章详情
	beego.Router("/Article/showArticleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/Article/UpdateArticle",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddArticleType;post:HandleAddType")
	//退出
	beego.Router("/Article/Logout",&controllers.ArticleController{},"get:Logout")
}
var FiltFunc = func(ctx *context.Context) {
	userName:=ctx.Input.Session("userName")
	if userName==nil {
		ctx.Redirect(302,"/")
	}
}