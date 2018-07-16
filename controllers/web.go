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
	"os"
	"strconv"
	"strings"
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
		if password == realpassword.Password {
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
func reqcert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		userName := sess.Get("username")
		notAfter, _ := strconv.Atoi(r.FormValue("notAfter"))
		ipAddress := strings.Split(r.FormValue("ipAddress"), ",")
		country := strings.Split(r.FormValue("country"), ",")
		organization := strings.Split(r.FormValue("organization"), ",")
		commonName := r.FormValue("commonName")
		ipPath := r.FormValue("ipAddress")
		//fmt.Println("chenyao*************", ipPath, ipAddress)
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("pubKeyFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		utils.Creatdir("conf/" + ipPath)
		f, err := os.OpenFile("conf/"+ipPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		utils.GetPubKey("conf/"+ipPath+"/key_req.pem", r.FormValue("ipAddress"))

		utils.GenerateCert(notAfter, ipAddress, country, organization, commonName)
		io.WriteString(w, "conf/"+ipPath+"/cert.pem")
		//fmt.Println("chenyao******", notAfter, ipAddress, country, organization, commonName)
		models.InsertData("cert", []string{"conf/" + ipPath + "/cert.pem", "conf/" + ipPath + "/pubkey.pem", userName.(string)})
	}
}

func easygen(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		userName := sess.Get("username")
		ipAddress := strings.Split(r.FormValue("ipAddress"), ",")
		utils.EasyGen(ipAddress)
		f1, _ := os.Open("conf/" + r.FormValue("ipAddress"))
		defer f1.Close()
		f2, _ := os.Open("conf/readme.txt")
		defer f2.Close()
		var files = []*os.File{f1, f2}
		dest := "conf/" + r.FormValue("ipAddress") + ".zip"
		utils.Compress(files, dest)
		models.InsertData("cert", []string{"conf/" + r.FormValue("ipAddress") + "/cert.pem", "conf/" + r.FormValue("ipAddress") + "/pubkey.pem", userName.(string)})
	}
}

func accountlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountlistMap := models.Queryaccountlist()
		accountlistJson, _ := json.Marshal(accountlistMap)
		io.WriteString(w, string(accountlistJson))
	}
}

func addaccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.FormValue("accountId")
		accountPassword := r.FormValue("accountPassword")
		models.InsertData("account", []string{accountId, accountPassword, "11"})
	}
}

func removeaccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		models.DeleteDateById("account", id)
	}
}

func certlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		iplistMap := models.QueryData("cert")
		iplistJson, _ := json.Marshal(iplistMap)
		io.WriteString(w, string(iplistJson))
	}
}
func createcrl(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		certPath := strings.Split(r.FormValue("certPath"), "@")
		utils.RevokedCertificates(certPath)
	}
}
func RunWeb() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态文件路径
	http.Handle("/conf/", http.StripPrefix("/conf/", http.FileServer(http.Dir("conf"))))       //设置静态文件路径
	http.HandleFunc("/login", login)                                                           //设置访问的路由
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/index_ca", index_ca)
	http.HandleFunc("/index_peer", index_peer)

	//用户操作
	http.HandleFunc("/addip", addip)
	http.HandleFunc("/removeip", removeip)
	http.HandleFunc("/iplist", iplist)
	http.HandleFunc("/reqcert", reqcert)
	http.HandleFunc("/easygen", easygen)

	//ca操作
	http.HandleFunc("/accountlist", accountlist)
	http.HandleFunc("/addaccount", addaccount)
	http.HandleFunc("/removeaccount", removeaccount)
	http.HandleFunc("/certlist", certlist)
	http.HandleFunc("/createcrl", createcrl)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
