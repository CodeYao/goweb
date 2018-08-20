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
		//globalSessions.SessionDestroy(w, r)

		sess.Delete("username")
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else {
		//定义session
		sess := globalSessions.SessionStart(w, r)
		accountId := r.FormValue("username")
		password := r.FormValue("password")
		password = utils.GetMD5(password)
		//请求的是登陆数据，那么执行登陆的逻辑判断
		realpassword := models.QueryPasswordByAccount(accountId)
		//fmt.Println(realpassword)

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
	//fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		userName := sess.Get("username")
		if userName == nil {
			t, _ := template.ParseFiles("views/login.html")
			t.Execute(w, nil)
		} else {
			t, _ := template.ParseFiles("views/index_ca.html")
			t.Execute(w, userName)
		}

	}
}
func index_peer(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	//fmt.Println("method:", r.Method) //获取请求的方法
	r.ParseForm()
	if r.Method == "GET" {
		userName := sess.Get("username")
		if userName == nil {
			t, _ := template.ParseFiles("views/login.html")
			t.Execute(w, nil)
		} else {
			t, _ := template.ParseFiles("views/index_peer.html")
			t.Execute(w, userName)
		}

	}
}
func showcertInfo(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			var certVO models.CertVO
			certid := r.FormValue("id")
			certInfo, _ := models.QueryDataById("cert", certid)
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
}
func certlist(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		var pagevo models.PageVO
		var startPage int
		var dataNum int
		if r.Method == "POST" {
			//userName := sess.Get("username")
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}

			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			pagevo.EntityList, _ = models.QuerycertlistByPage(userName.(string), strconv.Itoa(start), pagevo.PageNum)
			dataNum = models.GetTableNum("cert where accountId = '" + userName.(string) + "'")
			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}
			pagevolistJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pagevolistJson))
		}
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
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			//userName := sess.Get("username")
			//KeyType := r.FormValue("KeyType")
			KeyType := "sm2"
			certName := r.FormValue("certName")
			notAfter, _ := strconv.Atoi(r.FormValue("notAfter"))
			ipAddress := strings.Split(r.FormValue("ipAddress"), ";")
			reqcertxt := r.FormValue("reqcertxt")
			country := strings.Split(r.FormValue("country"), ";")
			organization := strings.Split(r.FormValue("organization"), ";")
			commonName := r.FormValue("commonName")
			//ipPath := r.FormValue("ipAddress")
			//fmt.Println(userName, certName, notAfter, ipAddress, reqcertxt, country, organization, commonName, ipPath, KeyType)
			var newpubkey string
			if KeyType == "sm2" {
				newpubkey = string(utils.GetSm2PubKey(reqcertxt))
			} else if KeyType == "ecdsa" {
				newpubkey = string(utils.GetEcdsaPubKey(reqcertxt))
			}
			//fmt.Println("************chenyao**********", newpubkey)
			if newpubkey == "" {
				//fmt.Println("************chenyao**********")
				io.WriteString(w, "1")
			} else {
				dataNum := models.GetTableNum("cert where certname = '" + certName + "'")
				if dataNum == 0 {
					newcert := string(utils.GenerateCert(notAfter, ipAddress, country, organization, commonName, reqcertxt, KeyType))
					//fmt.Println("newpubkey:", newpubkey)
					//fmt.Println("newcert:", newcert)
					models.InsertData("cert", []string{certName, r.FormValue("ipAddress"), newcert, newpubkey, userName.(string), reqcertxt, time.Now().Format("2006-01-02 15:04:05"), "待审批", "", "", "", ""})
				} else {
					io.WriteString(w, "0")
				}
			}

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
	var pagevo models.PageVO
	var startPage int
	var dataNum int
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			//selected := r.FormValue("selected")
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}

			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			//fmt.Println("chenyao*************************", pagevo.CurrentPage, pagevo.TotalPage)
			pagevo.EntityList, _ = models.QueryData("account where accountLevel = '11' and enabled = 'enabled' limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			dataNum = models.GetTableNum("account where accountLevel = '11' and enabled = 'enabled'")

			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}

			pageJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pageJson))
		}
	}

}

func addaccount(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := r.FormValue("accountId")
			dataNum := models.GetTableNum("account where accountId = '" + accountId + "' and enabled = 'enabled'")
			accountPassword := r.FormValue("accountPassword")
			if dataNum > 0 {
				io.WriteString(w, "err1")
			} else if len(accountId) > 8 {
				io.WriteString(w, "err2")
			} else if len(accountPassword) < 6 || len(accountPassword) > 12 {
				io.WriteString(w, "err3")
			} else {
				organization := r.FormValue("organization")
				accountPassword = utils.GetMD5(accountPassword)
				models.InsertData("account", []string{accountId, organization, accountPassword, "11", "enabled"})
			}

		}
	}

}

func removeaccount(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			//fmt.Println("chenyao**********", id)
			//models.DeleteDateById("account", id)
			models.DeleteAccountById(id)
			accountInfo, _ := models.QueryDataById("account", r.FormValue("id"))
			var sessionId string
			sessionMG := globalSessions.GetProvide()
			for k, v := range sessionMG.GetSession() {
				account := v.Value.(*SessionStore).GetValue()["username"]
				//fmt.Println(k, account)
				//fmt.Println(accountInfo["accountId"])
				if account.(string) == accountInfo["accountId"] {
					//fmt.Println("同时关闭session")
					sessionId = k
				}
			}
			sessionMG.SessionDestroy(sessionId)
		}
	}

}

func ca_certlist(w http.ResponseWriter, r *http.Request) {
	var pagevo models.PageVO
	var startPage int
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5
			dataNum := models.GetTableNum("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId")
			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}

			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			//fmt.Println("chenyao*************************", pagevo.CurrentPage, pagevo.TotalPage)
			pagevo.EntityList, _ = models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			pageJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pageJson))
		}
	}

}
func createcrl(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			var certPath []string
			certlistMap, _ := models.QueryData("cert where state = '已吊销'")
			for _, v := range certlistMap {
				//fmt.Println(k, "*********chenyao***********")
				certPath = append(certPath, v["cert"])
			}
			utils.RevokedCertificates(certPath)
		}
	}

}
func genreatePage(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		t, _ := template.ParseFiles("views/login.html")
		t.Execute(w, nil)
	} else {
		t, _ := template.ParseFiles("views/generate.html")
		t.Execute(w, nil)
	}

}
func changepwd(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			oldpwd := r.FormValue("oldpwd")
			oldpwd = utils.GetMD5(oldpwd)
			newpwd := r.FormValue("newpwd")
			newpwd = utils.GetMD5(newpwd)
			//请求的是登陆数据，那么执行登陆的逻辑判断
			realpassword := models.QueryPasswordByAccount(accountId.(string))
			//fmt.Println(realpassword)
			if oldpwd == realpassword.Password {
				models.UpdatePassword(accountId.(string), newpwd)
				io.WriteString(w, "修改成功，请重新登录")
			} else {
				io.WriteString(w, "旧密码不正确")
			}
		}
	}

}
func onetouch(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			KeyType := r.FormValue("KeyType")
			keysJson, _ := json.Marshal(append(utils.OneTouch(KeyType), accountId.(string)))
			io.WriteString(w, string(keysJson))
		}
	}

}
func aprovecert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			aprover := sess.Get("username")
			id := r.FormValue("id")
			dataNum := models.GetTableNum("ca where enabled = 'enabled'")
			if dataNum == 0 {
				io.WriteString(w, strconv.Itoa(dataNum))
			} else {
				models.UpdateCertAprove("已审批", aprover.(string), time.Now().Format("2006-01-02 15:04:05"), id)
			}

		}
	}
}
func rejectcert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			aprover := sess.Get("username")
			id := r.FormValue("id")
			models.UpdateCertAprove("已驳回", aprover.(string), time.Now().Format("2006-01-02 15:04:05"), id)
		}
	}

}
func revokecert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			revoker := sess.Get("username")
			id := r.FormValue("id")
			models.UpdateCertRevoke("已吊销", revoker.(string), time.Now().Format("2006-01-02 15:04:05"), id)
		}
	}

}
func select_certlist(w http.ResponseWriter, r *http.Request) {
	var pagevo models.PageVO
	var startPage int
	var dataNum int
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			selected := r.FormValue("selected")
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}
			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			//fmt.Println("chenyao*************************", pagevo.CurrentPage, pagevo.TotalPage)
			pagevo.EntityList, _ = models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			if selected == "全部" {
				dataNum = models.GetTableNum("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId")
				pagevo.EntityList, _ = models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			} else {
				dataNum = models.GetTableNum("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId where state = '" + selected + "'")
				pagevo.EntityList, _ = models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId where state = '" + selected + "' limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			}

			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}

			pageJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pageJson))
		}
	}

}

// func select_certlist(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		selected := r.FormValue("selected")
// 		if selected == "全部" {
// 			certlistMap := models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId")
// 			certlistJson, _ := json.Marshal(certlistMap)
// 			io.WriteString(w, string(certlistJson))
// 		} else {
// 			certlistMap := models.QueryData("(select * from tjfoc_ca.account where enabled = 'enabled') as A right join tjfoc_ca.cert as B on A.accountId = B.accountId where state = '" + selected + "'")
// 			certlistJson, _ := json.Marshal(certlistMap)
// 			io.WriteString(w, string(certlistJson))
// 		}

// 	}
// }
func downloadcert(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	accountId := sess.Get("username")
	if accountId == nil {
		io.WriteString(w, "<script>alert('登录超时，请重新登录'); window.location.href = 'login'</script>")
	} else {
		fileName := strconv.FormatInt(time.Now().Unix(), 10)
		certId := r.FormValue("certId")
		certstr, certName := models.QuerycertById(certId)
		fileName = certName + "_" + fileName
		//utils.ZipByte([]byte(certstr), "conf/cert.zip")
		ioutil.WriteFile("conf/"+fileName+".pem", []byte(certstr), 0666)
		f1, _ := os.Open("conf/" + fileName + ".pem")
		defer f1.Close()
		var files = []*os.File{f1}
		utils.Compress(files, "conf/"+fileName+".zip")
		//fmt.Println("chenyao***************", certId, certstr)
		file, _ := ioutil.ReadFile("conf/" + fileName + ".zip")
		w.Write(file)
		err := os.Remove("conf/" + fileName + ".pem")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func changepwd_ca(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := r.FormValue("cgaccountId")
			newpwd := r.FormValue("newpwd")
			newpwd = utils.GetMD5(newpwd)
			models.UpdatePassword(accountId, newpwd)
			io.WriteString(w, "修改成功")

		}
	}

}
func addaddress(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			address := r.FormValue("address")
			descript := r.FormValue("descript")
			models.InsertData("address", []string{address, time.Now().Format("2006-01-02 15:04:05"), accountId.(string), descript, "enabled"})
			io.WriteString(w, "添加成功")

		}
	}
}
func addresslist(w http.ResponseWriter, r *http.Request) {
	var pagevo models.PageVO
	var startPage int
	var dataNum int
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			//selected := r.FormValue("selected")
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}

			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			//fmt.Println("chenyao*************************", pagevo.CurrentPage, pagevo.TotalPage)
			pagevo.EntityList, _ = models.QueryData("address where enabled = 'enabled' limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			dataNum = models.GetTableNum("address where enabled = 'enabled'")

			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}

			pageJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pageJson))
		}
	}

}

func removeadderss(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			//fmt.Println("chenyao**********", id)
			models.DeleteAddressById(id)
		}
	}
}

func checkca(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			dataNum := models.GetTableNum("ca where enabled = 'enabled'")
			if dataNum == 0 {
				io.WriteString(w, strconv.Itoa(dataNum))
			} else {
				CAInfo, _ := models.QueryData("ca where enabled = 'enabled'")
				CAJson, _ := json.Marshal(CAInfo)
				io.WriteString(w, string(CAJson))
			}

		}
	}
}
func generateCA(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			//	keyType := r.FormValue("KeyType")
			keyType := "sm2"
			//caNotAfter := r.FormValue("canotAfter")
			caNotAfter, _ := strconv.Atoi(r.FormValue("canotAfter"))
			//caIPAddress := r.FormValue("caipAddress")
			caIPAddress := strings.Split(r.FormValue("caipAddress"), ";")
			//caCountry := r.FormValue("cacountry")
			caCountry := strings.Split(r.FormValue("cacountry"), ";")
			//caOrganization := r.FormValue("caorganization")
			caOrganization := strings.Split(r.FormValue("caorganization"), ";")
			caCommonName := r.FormValue("cacommonName")
			//fmt.Println("***chenyao***:", accountId, keyType, caNotAfter, caIPAddress, caCountry, caOrganization, caCommonName)
			caprivKey, capubKey := utils.GetCAKey(keyType)
			//fmt.Println(caprivKey, capubKey)
			newcert := string(utils.GenerateCACert(caNotAfter, caIPAddress, caCountry, caOrganization, caCommonName, caprivKey, keyType))
			//utils.GenerateCACert(keyType)
			//fmt.Println("capubKey:", capubKey)
			//fmt.Println("caprivKey:", caprivKey)
			//fmt.Println("newcert:", newcert)
			models.InsertData("ca", []string{capubKey, caprivKey, newcert, keyType, time.Now().Format("2006-01-02 15:04:05"), accountId.(string), "enabled"})
		}
	}

}

func createcacert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cacert, _ := models.QueryData("ca where enabled = 'enabled'")
		ioutil.WriteFile("conf/cacert.pem", []byte(cacert[0]["cacert"]), 0666)
	}
}
func importCA(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			//	keyType := r.FormValue("KeyType")
			keyType := "sm2"
			caprivkey := r.FormValue("caprivkeytext")
			cacert := r.FormValue("cacertxt")
			capubkey, err := utils.GetSm2PubKeyFromPrivKey(caprivkey)
			if !(utils.CheckPrivAndCert([]byte(cacert), capubkey)) {
				io.WriteString(w, "err0")
			} else if err != nil {
				//fmt.Println(err, caprivkey)
				io.WriteString(w, "err1")
			} else {
				err = utils.VerifySM2(cacert)
				if err != nil {
					//fmt.Println(err)
					io.WriteString(w, "err2")
				} else {
					models.InsertData("ca", []string{string(capubkey), caprivkey, cacert, keyType, time.Now().Format("2006-01-02 15:04:05"), accountId.(string), "enabled"})
				}
			}

		}
	}
}

func codeaddresslist(w http.ResponseWriter, r *http.Request) {
	var pagevo models.PageVO
	var startPage int
	var dataNum int
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			//selected := r.FormValue("selected")
			pagevo.CurrentPage = r.FormValue("currentpage")
			if pagevo.CurrentPage == "" {
				pagevo.CurrentPage = "1"
			}
			currentPage, _ := strconv.Atoi(pagevo.CurrentPage)
			pagevo.PageNum = "5"
			pageNum := 5
			showPage := 5

			start := (currentPage - 1) * pageNum
			if currentPage%showPage == 0 {
				startPage = (currentPage - showPage) + 1
			} else {
				startPage = currentPage/showPage*showPage + 1
			}

			//fmt.Println(startPage, "chenyao********************", currentPage)
			pagevo.StartPage = strconv.Itoa(startPage)
			//fmt.Println("chenyao*************************", pagevo.CurrentPage, pagevo.TotalPage)
			pagevo.EntityList, _ = models.QueryData("codeaddress where enabled = 'enabled' limit " + strconv.Itoa(start) + "," + pagevo.PageNum)
			dataNum = models.GetTableNum("codeaddress where enabled = 'enabled'")

			if dataNum%pageNum == 0 {
				pagevo.TotalPage = strconv.Itoa(dataNum / pageNum)
			} else {
				pagevo.TotalPage = strconv.Itoa((dataNum / pageNum) + 1)
			}

			pageJson, _ := json.Marshal(pagevo)
			io.WriteString(w, string(pageJson))
		}
	}

}
func codeaddaddress(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			accountId := sess.Get("username")
			address := r.FormValue("address")
			descript := r.FormValue("descript")
			models.InsertData("codeaddress", []string{address, time.Now().Format("2006-01-02 15:04:05"), accountId.(string), descript, "enabled"})
			io.WriteString(w, "添加成功")

		}
	}

}

func removecodeadderss(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	userName := sess.Get("username")
	if userName == nil {
		io.WriteString(w, "Timeout")
	} else {
		if r.Method == "POST" {
			id, _ := strconv.Atoi(r.FormValue("id"))
			//fmt.Println("chenyao**********", id)
			models.DeletecodeAddressById(id)
		}
	}

}

func RunWeb() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //设置静态文件路径
	http.Handle("/conf/", http.StripPrefix("/conf/", http.FileServer(http.Dir("conf"))))       //设置静态文件路径

	http.HandleFunc("/login", login) //设置访问的路由
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
	http.HandleFunc("/changepwd_ca", changepwd_ca)
	http.HandleFunc("/addaddress", addaddress)
	http.HandleFunc("/addresslist", addresslist)
	http.HandleFunc("/removeadderss", removeadderss)
	http.HandleFunc("/checkca", checkca)
	http.HandleFunc("/generateCA", generateCA)
	http.HandleFunc("/createcacert", createcacert)
	http.HandleFunc("/importCA", importCA)
	http.HandleFunc("/codeaddresslist", codeaddresslist)
	http.HandleFunc("/codeaddaddress", codeaddaddress)
	http.HandleFunc("/removecodeadderss", removecodeadderss)

	// //配置rpc方法
	// var address = new(carpc.Address)
	// rpc.Register(address)
	// rpc.HandleHTTP() //将Rpc绑定到HTTP协议上。
	fmt.Println("启动服务...输入crtl+c退出服务")
	err := http.ListenAndServe(":9093", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
