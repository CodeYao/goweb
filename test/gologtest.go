package main

import (
	"ca/goweb/controllers"
	"fmt"
	"log"
	"os"
)

func test_deferfatal() {
	defer func() {
		fmt.Println("--first--")
	}()
	log.Fatalln("test for defer Fatal")
}

func main() {
	arr := []int{2, 3}
	log.Print("Print array ", arr, "\n")
	log.Println("Println array", arr)
	log.Printf("Printf array with item [%d,%d]\n", arr[0], arr[1])
	//test_deferfatal()

	fileName := "Info_First.log"
	logFile, err := os.Create(fileName)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error")
	}
	debugLog := log.New(logFile, "[Info]", log.Llongfile)
	debugLog.Println("A Info message here")
	debugLog.SetPrefix("[Debug]")
	debugLog.Println("A Debug Message here ")
	controllers.Config("./logs.log", controllers.InfoLevel)
	controllers.Debugf("hahahahh")
	controllers.Infof("xx", "def", "ghijk")
	controllers.Warnf("xxxx", "def", "ghijk")
	controllers.Errorf("xxxxx", "def", "ghijk")
}
