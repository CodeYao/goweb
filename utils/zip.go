package utils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipByte(certByte []byte, path string) {
	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"cert.pem", string(certByte)},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	buf.WriteTo(f)
}

func GetMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

// func main() {
// 	f1, _ := os.Open("../conf/10.1.1.2")

// 	defer f1.Close()
// 	f2, _ := os.Open("../conf/ca")

// 	defer f2.Close()
// 	f3, _ := os.Open("../conf/10.1.1.2/cert.pem")

// 	defer f3.Close()
// 	var files = []*os.File{f1, f2, f3}
// 	dest := "test.zip"
// 	Compress(files, dest)

// }
