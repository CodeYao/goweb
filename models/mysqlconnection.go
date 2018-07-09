package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbhostsip  = "127.0.0.1"
	dbusername = "root"
	dbpassowrd = "root"
	dbname     = "tjfoc"
)

type mysql_db struct {
	db *sql.DB //定义结构体
}

func (f *mysql_db) mysql_open() { //打开
	Odb, err := sql.Open("mysql", dbusername+":"+dbpassowrd+"@tcp("+dbhostsip+")/"+dbname)
	if err != nil {
		fmt.Println("链接失败")
	}
	fmt.Println("链接数据库成功...........已经打开")
	f.db = Odb
}

func (f *mysql_db) mysql_close() { //关闭
	defer f.db.Close()
	fmt.Println("链接数据库成功...........已经关闭")
}

func (f *mysql_db) mysql_select(sql_data string) {
	rows, err := f.db.Query(sql_data)
	if err != nil {
		println(err)
	}
	for rows.Next() {
		var in_param string

		err = rows.Scan(&in_param)
		if err != nil {
			panic(err)
		}
	}
}
