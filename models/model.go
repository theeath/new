package models

import ("github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string
	Passwd string
	Articles[]*Article`orm:"rel(m2m)"`
}
type Article struct {
	Id int`orm:"pk;auto"`
	Title string`orm:"size(20)"`
	Content string`orm:"size(500)"`
	Img string		`orm:"size(200);null"`
	Time time.Time`orm:"type(datetime);auto_now_add"`
	Count int`orm:"default(0)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users[]*User`orm:"reverse(many)"`
}
type ArticleType struct {
	Id int
	TypeName string `orm:"size(40)"`
	Articles[]*Article`orm:"reverse(many)"`
}
func init()  {
	orm.RegisterDataBase("default","mysql","root:960825fa@tcp(127.0.0.1:3306)/newsWeb?charset=utf8&loc=Local")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}