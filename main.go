package main

import (
	"github.com/astaxie/beego"
	_ "news/models"
	_ "news/routers"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowNextPage",HandleNextPage)
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.Run()
}
func HandlePrePage(data int)(string)  {
	pageIndex := data - 1

	pageIndex1:=strconv.Itoa(pageIndex)
	return pageIndex1

}
func HandleNextPage(data int)(string)  {

	pageIndex := data + 1
	pageIndex1:=strconv.Itoa(pageIndex)
	return pageIndex1
}

