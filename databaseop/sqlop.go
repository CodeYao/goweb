package databaseop

import (
	"fmt"
	"strings"
)

var db mysql_db

func TestSelect() {
	db.mysql_open()
	//查询数据，取所有字段
	rows2, _ := db.db.Query("select * from t_test")
	db.mysql_close()
	//返回所有列
	cols, _ := rows2.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]string)
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	fmt.Println(result)
}

func insertData(tableName string, insertdata []string) {
	db.mysql_open()
	sql := "insert into " + tableName + " values(null," + strings.Join(insertdata, ",") + ")"
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
	row := db.db.QueryRow(sql, 1)
	db.mysql_close()
	//result := make(map[int]map[string]string)
	fmt.Println(row)
}
