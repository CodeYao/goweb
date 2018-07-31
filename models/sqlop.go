package models

import (
	"fmt"
	"strings"
)

var db mysql_db

func Queryaccountlist() []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from account where accountLevel = ? and enabled = 'enabled'", "11")
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result []map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result, row)
		i++
	}
	return result
}
func QuerycertlistByPage(accountId string, start string, pageNum string) []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from cert where accountId = ? limit ?,?", accountId, start, pageNum)
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result []map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result, row)
		i++
	}
	return result
}
func QuerycertlistByAccountId(accountId string) []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from cert where accountId = ?", accountId)
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result []map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result, row)
		i++
	}
	return result
}

func QueryData(tableName string) []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from " + tableName)
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result []map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result, row)
		i++
	}
	return result
}

func DeleteDateById(tableName string, id string) {
	db.mysql_open()
	ret, _ := db.db.Exec("delete from "+tableName+" where id = ?", id)
	//获取影响行数
	db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	fmt.Println(del_nums)
}

//将账号更新为不可用
func DeleteAccountById(id int) {
	db.mysql_open()
	ret, _ := db.db.Exec("update account set enabled = 'disabled' where Id = ?", id)
	db.mysql_close()
	fmt.Println("删除成功", ret.RowsAffected)
}

//将钱包地址更新为不可用
func DeleteAddressById(id int) {
	db.mysql_open()
	ret, _ := db.db.Exec("update address set enabled = 'disabled' where Id = ?", id)
	db.mysql_close()
	fmt.Println("删除成功", ret.RowsAffected)
}

//将冻结钱包地址更新为不可用
func DeletecodeAddressById(id int) {
	db.mysql_open()
	ret, _ := db.db.Exec("update codeaddress set enabled = 'disabled' where Id = ?", id)
	db.mysql_close()
	fmt.Println("删除成功", ret.RowsAffected)
}

func InsertData(tableName string, insertdata []string) {
	db.mysql_open()
	sql := "insert into " + tableName + " values(null,'" + strings.Join(insertdata, "','") + "')"
	fmt.Println(sql)
	ret, err := db.db.Exec(sql)
	db.mysql_close()
	if err != nil {
		panic(err)
	}
	//获取插入ID
	fmt.Println("插入成功", ret, sql)
}

func QueryDataById(tableName string, id string) map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from "+tableName+" where id =?", id)
	if err != nil {
		panic(err)
	}
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = row
		i++
	}
	return result
}

//根据账号查询密码
func QueryPasswordByAccount(accountId string) Account {
	var account Account
	db.mysql_open()
	row := db.db.QueryRow("select accountId,organization,password,accountLevel,enabled from account where accountId = '" + accountId + "'")
	db.mysql_close()
	row.Scan(&account.AccountId, &account.Organization, &account.Password, &account.AccountLevel, &account.Enable)
	return account
}

//修改账号密码
func UpdatePassword(accountId string, password string) {
	db.mysql_open()
	ret, _ := db.db.Exec("update account set password = ? where accountId = ?", password, accountId)
	db.mysql_close()
	fmt.Println("修改成功", ret)
}

//修改审批状态
func UpdateCertAprove(state string, aprover string, aproveDate string, certid string) {
	db.mysql_open()
	ret, _ := db.db.Exec("update cert set state =? , approver = ? , approvaldate = ? where id = ?", state, aprover, aproveDate, certid)
	db.mysql_close()
	fmt.Println("修改状态成功", ret)
}

//修改审批状态
func UpdateCertRevoke(state string, revoker string, revokeDate string, certid string) {
	db.mysql_open()
	ret, _ := db.db.Exec("update cert set state =? , revoker = ? , revokedate = ? where id = ?", state, revoker, revokeDate, certid)
	db.mysql_close()
	fmt.Println("吊销成功", ret)
}

//获取证书字符串
func QuerycertById(id string) string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from cert where id =?", id)
	if err != nil {
		panic(err)
	}
	db.mysql_close()
	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	var result map[string]string
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = row
		i++
	}
	return result["cert"]
}

func GetTableNum(table string) int {
	var dataCount int
	db.mysql_open()
	//查询数据，取所有字段
	row := db.db.QueryRow("select count(-1) from " + table)
	db.mysql_close()
	err := row.Scan(&dataCount)
	if err != nil {
		panic(err)
	}
	// return strconv.Itoa(dataCount)
	return dataCount
}
