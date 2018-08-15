package utils

import (
	"io/ioutil"
	"os"
)

func CheckIsExist(fileName string) bool {
	var exist = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func write(filename string, b string) {
	var err1 error

	if !checkFileIsExist(filename) { //如果文件存在
		_, err1 = os.Create(filename) //创建文件
	}
	if err1 != nil {
		panic(err1)
		return
	}

	context, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(context))

	context = append(context, []byte(b)...)

	ioutil.WriteFile(filename, context, os.ModeAppend)

}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
