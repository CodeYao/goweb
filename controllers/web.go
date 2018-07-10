package controllers

import (
	"ca/goweb/models"
	"ca/goweb/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

var globalSessions *utils.Manager

//初始化session
func init() {
	globalSessions, _ = utils.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}
func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sess := globalSessions.SessionStart(w, r)
		sess.Delete("username")
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else {
		//定义session
		sess := globalSessions.SessionStart(w, r)

		accountId := r.FormValue("username")
		password := r.FormValue("password")
		//请求的是登陆数据，那么执行登陆的逻辑判断
		realpassword := models.QueryPasswordByAccount(accountId)
		fmt.Println(realpassword)
		if password == realpassword.AccountPassword {
			//str, _ := json.Marshal(realpassword)
			//设置session
			sess.Set("username", accountId)
			//http.Redirect(w, r, "/", 302)
			io.WriteString(w, realpassword.AccountLevel)
		} else {
			io.WriteString(w, "账号或者密码错误")
		}
	}
}
func index_ca(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		userName := sess.Get("username")
		t, _ := template.ParseFiles("views/index_ca.html")
		t.Execute(w, userName)
	}
}
func index_peer(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		userName := sess.Get("username")
		t, _ := template.ParseFiles("views/index_peer.html")
		t.Execute(w, userName)
	}
}
func addip(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		userName := sess.Get("username")
		ipname := r.FormValue("ipname")
		ipstr := r.FormValue("ipstr")
		models.InsertData("iplist", []string{ipname, ipstr, userName.(string)})
	}
}
func iplist(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		userName := sess.Get("username")
		iplistMap := models.QueryiplistByAccountId(userName.(string))
		iplistJson, _ := json.Marshal(iplistMap)
		io.WriteString(w, string(iplistJson))
	}
}
func removeip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ipid := r.FormValue("id")
		models.DeleteDateById("iplist", ipid)
	}
}
func RunWeb() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态文件路径
	http.HandleFunc("/login", login)                                                           //设置访问的路由
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/index_ca", index_ca)
	http.HandleFunc("/index_peer", index_peer)
	http.HandleFunc("/addip", addip)
	http.HandleFunc("/removeip", removeip)
	http.HandleFunc("/iplist", iplist)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
