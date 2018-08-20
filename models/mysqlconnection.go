package models

import (
	"database/sql"
	"encoding/xml"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// var (
// 	dbhostsip  = "127.0.0.1"
// 	dbusername = "root"
// 	dbpassowrd = "root"
// 	dbname     = "tjfoc_ca"
// )

type mysql_db struct {
	db *sql.DB //定义结构体
}

func newDBEngine(f string) (*mysql_db, error) {
	// file, err := os.Open("static/dataconf.xml")
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	d := xml.NewDecoder(file)
	var resource Resource
	err = d.Decode(&resource)
	if err != nil {
		return nil, err
	}
	eng := &mysql_db{}
	eng.db, err = sql.Open("mysql", resource.Dbusername+":"+resource.Dbpassowrd+"@tcp("+resource.Dbhostsip+")/"+resource.Dbname)
	if err != nil {
		return nil, err
	}
	err = eng.db.Ping()
	if err != nil {
		return nil, err
	}
	return eng, nil
}

func (f *mysql_db) mysql_open() { //打开
	// file, err := os.Open("static/dataconf.xml")

	// if err != nil {
	// 	//panic(err)
	// 	file, err = os.Open("../grpcdataconf/dataconf.xml")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fileinfo, _ := file.Stat()
	// filesize := fileinfo.Size()
	// buffer := make([]byte, filesize)
	// file.Read(buffer)
	// defer file.Close()
	// //fmt.Printf("%s", buffer)
	// xml.Unmarshal(buffer, &resource)
	// Odb, err := sql.Open("mysql", resource.Dbusername+":"+resource.Dbpassowrd+"@tcp("+resource.Dbhostsip+")/"+resource.Dbname)
	// if err != nil {
	// 	fmt.Println("数据库链接失败", err)
	// }
	// //fmt.Println("链接数据库成功...........已经打开")
	// f.db = Odb

}

func (f *mysql_db) mysql_close() { //关闭
	if f.db != nil {
		f.db.Close()
	}
	//fmt.Println("链接数据库成功...........已经关闭")
}

func (f *mysql_db) mysql_select(sql_data string) error {
	rows, err := f.db.Query(sql_data)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var in_param string

		err = rows.Scan(&in_param)
		if err != nil {
			return err
		}
	}
	return nil
}
