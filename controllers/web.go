package controllers

import (
	"ca/goweb/models"
	"ca/goweb/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
		if password == realpassword.Password && realpassword.Enable == "enabled" {
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
func showcertInfo(w http.ResponseWriter, r *http.Request) {
	//sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		var certVO models.CertVO
		//userName := sess.Get("username")
		certid := r.FormValue("id")
		certInfo := models.QueryDataById("cert", certid)
		certstr := certInfo["cert"]
		cert := utils.GetCert([]byte(certstr))
		certVO.CertName = certInfo["certname"]
		certVO.IpAdderss = certInfo["ipstr"] //strings.Join(cert.DNSNames, ";")
		certVO.CertDay = cert.NotAfter.Format("2006-01-02")
		certVO.Country = strings.Join(cert.Subject.Country, ";")
		certVO.Organization = strings.Join(cert.Subject.Organization, ";")
		certVO.CommonName = cert.Subject.CommonName
		certVO.State = certInfo["state"]
		certInfoJson, _ := json.Marshal(certVO)
		io.WriteString(w, string(certInfoJson))
	}
}
func certlist(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		userName := sess.Get("username")
		certlistMap := models.QuerycertlistByAccountId(userName.(string))
		certlistJson, _ := json.Marshal(certlistMap)
		io.WriteString(w, string(certlistJson))
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
		certName := r.FormValue("certName")
		notAfter, _ := strconv.Atoi(r.FormValue("notAfter"))
		ipAddress := strings.Split(r.FormValue("ipAddress"), ";")
		reqcertxt := r.FormValue("reqcertxt")
		country := strings.Split(r.FormValue("country"), ";")
		organization := strings.Split(r.FormValue("organization"), ";")
		commonName := r.FormValue("commonName")
		ipPath := r.FormValue("ipAddress")
		fmt.Println(userName, certName, notAfter, ipAddress, reqcertxt, country, organization, commonName, ipPath)
		newpubkey := string(utils.GetPubKey(reqcertxt))
		newcert := string(utils.GenerateCert(notAfter, ipAddress, country, organization, commonName, reqcertxt))
		fmt.Println("newpubkey:", newpubkey)
		fmt.Println("newcert:", newcert)
		models.InsertData("cert", []string{certName, r.FormValue("ipAddress"), newcert, newpubkey, userName.(string), time.Now().Format("2006-01-02 15:04:05"), "待审批", "", "", "", ""})
		//fmt.Println("chenyao*************", ipPath, ipAddress)
		// r.ParseMultipartForm(32 << 20)
		// file, handler, err := r.FormFile("pubKeyFile")
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// defer file.Close()
		// //fmt.Fprintf(w, "%v", handler.Header)
		// utils.Creatdir("conf/" + ipPath)
		// f, err := os.OpenFile("conf/"+ipPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// defer f.Close()
		// io.Copy(f, file)

		// utils.GetPubKey("conf/"+ipPath+"/key_req.pem", r.FormValue("ipAddress"))

		// utils.GenerateCert(notAfter, ipAddress, country, organization, commonName)
		// io.WriteString(w, "conf/"+ipPath+"/cert.pem")
		// //fmt.Println("chenyao******", notAfter, ipAddress, country, organization, commonName)
		// models.InsertData("cert", []string{"conf/" + ipPath + "/cert.pem", "conf/" + ipPath + "/pubkey.pem", userName.(string)})
	}
}

func easygen(w http.ResponseWriter, r *http.Request) {
	// sess := globalSessions.SessionStart(w, r)
	// if r.Method == "POST" {
	// 	userName := sess.Get("username")
	// 	ipAddress := strings.Split(r.FormValue("ipAddress"), ";")
	// 	utils.EasyGen(ipAddress)
	// 	f1, _ := os.Open("conf/" + r.FormValue("ipAddress"))
	// 	defer f1.Close()
	// 	f2, _ := os.Open("conf/readme.txt")
	// 	defer f2.Close()
	// 	var files = []*os.File{f1, f2}
	// 	dest := "conf/" + r.FormValue("ipAddress") + ".zip"
	// 	utils.Compress(files, dest)
	// 	models.InsertData("cert", []string{"conf/" + r.FormValue("ipAddress") + "/cert.pem", "conf/" + r.FormValue("ipAddress") + "/pubkey.pem", userName.(string)})
	// }
}

func accountlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("chenyao*******************")
		accountlistMap := models.Queryaccountlist()
		accountlistJson, _ := json.Marshal(accountlistMap)
		io.WriteString(w, string(accountlistJson))
	}
}

func addaccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		accountId := r.FormValue("accountId")
		accountPassword := r.FormValue("accountPassword")
		organization := r.FormValue("organization")
		models.InsertData("account", []string{accountId, organization, accountPassword, "11", "enabled"})
	}
}

func removeaccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		//models.DeleteDateById("account", id)
		models.DeleteAccountById(id)
	}
}

func ca_certlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		certlistMap := models.QueryData("tjfoc_ca.cert as A left join tjfoc_ca.account as B on A.accountId = B.accountId")
		certlistJson, _ := json.Marshal(certlistMap)
		io.WriteString(w, string(certlistJson))
	}
}
func createcrl(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var certPath []string
		certlistMap := models.QueryData("cert where state = '已吊销'")
		for k, v := range certlistMap {
			fmt.Println(k, "*********chenyao***********")
			certPath = append(certPath, v["cert"])
		}
		utils.RevokedCertificates(certPath)
	}
}
func genreatePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/generate.html")
	t.Execute(w, nil)
}
func changepwd(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		accountId := sess.Get("username")
		oldpwd := r.FormValue("oldpwd")
		newpwd := r.FormValue("newpwd")
		//请求的是登陆数据，那么执行登陆的逻辑判断
		realpassword := models.QueryPasswordByAccount(accountId.(string))
		fmt.Println(realpassword)
		if oldpwd == realpassword.Password {
			models.UpdatePassword(accountId.(string), newpwd)
			io.WriteString(w, "修改成功，请重新登录")
		} else {
			io.WriteString(w, "旧密码不正确")
		}
	}
}
func onetouch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		keysJson, _ := json.Marshal(utils.OneTouch())
		io.WriteString(w, string(keysJson))
	}
}
func aprovecert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		aprover := sess.Get("username")
		id := r.FormValue("id")
		models.UpdateCertAprove("已审批", aprover.(string), time.Now().Format("2006-01-02 15:04:05"), id)
	}
}
func rejectcert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		aprover := sess.Get("username")
		id := r.FormValue("id")
		models.UpdateCertAprove("已驳回", aprover.(string), time.Now().Format("2006-01-02 15:04:05"), id)
	}
}
func revokecert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	if r.Method == "POST" {
		revoker := sess.Get("username")
		id := r.FormValue("id")
		models.UpdateCertRevoke("已吊销", revoker.(string), time.Now().Format("2006-01-02 15:04:05"), id)
	}
}
func select_certlist(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		selected := r.FormValue("selected")
		if selected == "全部" {
			certlistMap := models.QueryData("tjfoc_ca.cert as A left join tjfoc_ca.account as B on A.accountId = B.accountId")
			certlistJson, _ := json.Marshal(certlistMap)
			io.WriteString(w, string(certlistJson))
		} else {
			certlistMap := models.QueryData("tjfoc_ca.cert as A left join tjfoc_ca.account as B on A.accountId = B.accountId where state = '" + selected + "'")
			certlistJson, _ := json.Marshal(certlistMap)
			io.WriteString(w, string(certlistJson))
		}

	}
}
func downloadcert(w http.ResponseWriter, r *http.Request) {
	fileName := strconv.FormatInt(time.Now().Unix(), 10)
	certId := r.FormValue("certId")
	certstr := models.QuerycertById(certId)
	//utils.ZipByte([]byte(certstr), "conf/cert.zip")
	ioutil.WriteFile("conf/"+fileName+".pem", []byte(certstr), 0666)
	f1, _ := os.Open("conf/" + fileName + ".pem")
	defer f1.Close()
	var files = []*os.File{f1}
	utils.Compress(files, "conf/"+fileName+".zip")
	fmt.Println("chenyao***************", certId, certstr)
	file, _ := ioutil.ReadFile("conf/" + fileName + ".zip")
	w.Write(file)
	err := os.Remove("conf/" + fileName + ".pem")
	if err != nil {
		fmt.Println(err)
	}
}
func RunWeb() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态文件路径
	http.Handle("/conf/", http.StripPrefix("/conf/", http.FileServer(http.Dir("conf"))))       //设置静态文件路径
	http.HandleFunc("/login", login)                                                           //设置访问的路由
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/changepwd", changepwd)
	http.HandleFunc("/index_ca", index_ca)
	http.HandleFunc("/index_peer", index_peer)

	//用户操作
	http.HandleFunc("/showcertInfo", showcertInfo)
	http.HandleFunc("/removeip", removeip)
	http.HandleFunc("/reqcert", reqcert)
	http.HandleFunc("/easygen", easygen)
	http.HandleFunc("/certlist", certlist)
	http.HandleFunc("/genreatePage", genreatePage)
	http.HandleFunc("/onetouch", onetouch)

	//ca操作
	http.HandleFunc("/accountlist", accountlist)
	http.HandleFunc("/addaccount", addaccount)
	http.HandleFunc("/removeaccount", removeaccount)
	//http.HandleFunc("/certlist", certlist)
	http.HandleFunc("/ca_certlist", ca_certlist)
	http.HandleFunc("/createcrl", createcrl)
	http.HandleFunc("/aprovecert", aprovecert)
	http.HandleFunc("/rejectcert", rejectcert)
	http.HandleFunc("/revokecert", revokecert)
	http.HandleFunc("/select_certlist", select_certlist)
	http.HandleFunc("/downloadcert", downloadcert)

	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
