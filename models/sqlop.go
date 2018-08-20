package models

import (
	"CAWeb/models"
	"fmt"
	"strings"
)

var db *mysql_db

func init() {
	var err error
	db, err = newDBEngine("static/dataconf.xml")
	if err != nil {
		Fatalf("database connection failed [%v]", err)
	}

}

func Queryaccountlist() ([]map[string]string, error) {
	// db.mysql_open()
	// defer db.mysql_close()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from account where accountLevel = ? and enabled = 'enabled'", "11")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//返回所有列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
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
	return result, nil
}
func QuerycertlistByPage(accountId string, start string, pageNum string) ([]map[string]string, error) {
	// db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from cert where accountId = ? limit ?,?", accountId, start, pageNum)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// db.mysql_close()
	//返回所有列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
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
	return result, nil
}
func QuerycertlistByAccountId(accountId string) ([]map[string]string, error) {
	// db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from cert where accountId = ?", accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// db.mysql_close()
	//返回所有列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
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
	return result, nil
}

func QueryData(tableName string) ([]map[string]string, error) {
	// db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// db.mysql_close()
	//返回所有列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
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
	return result, nil
}

func DeleteDateById(tableName string, id string) error {
	// db.mysql_open()
	ret, err := db.db.Exec("delete from "+tableName+" where id = ?", id)
	if err != nil {
		return err
	}

	//获取影响行数
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by delete success is %v", del_nums)
	//fmt.Println("删除成功，受影响行数为：", del_nums)
	return nil
}

//将账号更新为不可用
func DeleteAccountById(id int) error {
	// db.mysql_open()
	ret, err := db.db.Exec("update account set enabled = 'disabled' where Id = ?", id)
	if err != nil {
		return err
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by delete success is %v", del_nums)
	//fmt.Println("删除成功，受影响行数为：", del_nums)
	return nil
}

//将钱包地址更新为不可用
func DeleteAddressById(id int) error {
	// db.mysql_open()
	ret, err := db.db.Exec("update address set enabled = 'disabled' where Id = ?", id)
	if err != nil {
		return err
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by delete success is %v", del_nums)
	return nil
}

//将冻结钱包地址更新为不可用
func DeletecodeAddressById(id int) error {
	// db.mysql_open()
	ret, err := db.db.Exec("update codeaddress set enabled = 'disabled' where Id = ?", id)
	if err != nil {
		// panic(err)
		return err
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by delete success is %v", del_nums)
	//fmt.Println("删除成功，受影响行数为：", del_nums)
	return nil
}

func InsertData(tableName string, insertdata []string) error {
	// db.mysql_open()
	sql := "insert into " + tableName + " values(null,'" + strings.Join(insertdata, "','") + "')"
	//fmt.Println(sql)
	ret, err := db.db.Exec(sql)
	if err != nil {
		// panic(err)
		return nil
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by insertion success is %v", del_nums)
	//fmt.Println("插入成功，受影响行数为：", del_nums)
	//获取插入ID
	//fmt.Println("插入成功", ret, sql)
	return nil
}

func QueryDataById(tableName string, id string) (map[string]string, error) {
	// db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from "+tableName+" where id =?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// db.mysql_close()
	//返回所有列
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
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
	return result, nil
}

//根据账号查询密码
func QueryPasswordByAccount(accountId string) Account {
	var account Account
	// db.mysql_open()
	row := db.db.QueryRow("select accountId,organization,password,accountLevel,enabled from account where accountId = '" + accountId + "'")
	// db.mysql_close()
	err := row.Scan(&account.AccountId, &account.Organization, &account.Password, &account.AccountLevel, &account.Enable)
	if err != nil {
		models.Errorf("query pwd error: %s", err)
	}
	return account
}

//修改账号密码
func UpdatePassword(accountId string, password string) {
	// db.mysql_open()
	ret, err := db.db.Exec("update account set password = ? where accountId = ?", password, accountId)
	if err != nil {
		models.Errorf("update pwd error: %s", err)
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by update success is %v", del_nums)
	//fmt.Println("修改成功", ret)
}

//修改审批状态
func UpdateCertAprove(state string, aprover string, aproveDate string, certid string) {
	// db.mysql_open()
	ret, err := db.db.Exec("update cert set state =? , approver = ? , approvaldate = ? where id = ?", state, aprover, aproveDate, certid)
	if err != nil {
		models.Errorf("update approval status error: %s", err)
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by update success is %d", del_nums)
	//fmt.Println("修改状态成功", ret)
}

//修改审批状态
func UpdateCertRevoke(state string, revoker string, revokeDate string, certid string) {
	// db.mysql_open()
	ret, err := db.db.Exec("update cert set state =? , revoker = ? , revokedate = ? where id = ?", state, revoker, revokeDate, certid)
	if err != nil {

		fmt.Println("update approval status error:", err)
	}
	// db.mysql_close()
	del_nums, _ := ret.RowsAffected()
	models.Infof("The number of rows affected by insertion success is %d", del_nums)
	//fmt.Println("插入成功，受影响行数为：", del_nums)
	//fmt.Println("吊销成功", ret)
}

//获取证书字符串
func QuerycertById(id string) (string, string) {
	// db.mysql_open()
	//查询数据，取所有字段
	rows, err := db.db.Query("select * from cert where id =?", id)
	if err != nil {
		fmt.Println("Query cert by id error:", err)
		return "", ""
	}
	defer rows.Close()
	// db.mysql_close()
	//返回所有列
	cols, err := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	if err != nil {
		panic(err)
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
	return result["cert"], result["certname"]
}

func GetTableNum(table string) int {
	var dataCount int
	// db.mysql_open()
	//查询数据，取所有字段
	row := db.db.QueryRow("select count(1) from " + table)

	err := row.Scan(&dataCount)
	if err != nil {
		models.Errorf("get table num error:", err)
	}
	// db.mysql_close()
	// return strconv.Itoa(dataCount)
	return dataCount
}
