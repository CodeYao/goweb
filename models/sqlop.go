package models

import (
	"fmt"
	"strings"
)

var db mysql_db

func Queryaccountlist() []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from account where accountLevel = ?", "11")
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

func QueryiplistByAccountId(accountId string) []map[string]string {
	db.mysql_open()
	//查询数据，取所有字段
	rows, _ := db.db.Query("select * from iplist where accountId = ?", accountId)
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

func queryDataById(tableName string, id string) {
	db.mysql_open()
	sql := "select * from " + tableName + " where Id = " + id
	row, _ := db.db.Query(sql)
	db.mysql_close()
	//result := make(map[int]map[string]string)
	fmt.Println(row)
}

//根据账号查询密码
func QueryPasswordByAccount(accountId string) Account {
	var account Account
	db.mysql_open()
	row := db.db.QueryRow("select accountId,accountPassword,accountLevel from account where accountId = '" + accountId + "'")
	db.mysql_close()
	row.Scan(&account.AccountId, &account.AccountPassword, &account.AccountLevel)
	return account
}
