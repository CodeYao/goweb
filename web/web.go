package web

import (
	"ca/goweb/models"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else {
		accountId := r.FormValue("username")
		password := r.FormValue("password")
		//请求的是登陆数据，那么执行登陆的逻辑判断
		realpassword := models.QueryPasswordByAccount(accountId)
		fmt.Println(realpassword)
		if password == realpassword[0]["accountPassword"] {
			io.WriteString(w, "登录成功")
		} else {
			io.WriteString(w, "账号或者密码错误")
		}
	}
}

func RunWeb() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态文件路径
	http.HandleFunc("/login", login)                                                           //设置访问的路由
	err := http.ListenAndServe(":9090", nil)                                                   //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
