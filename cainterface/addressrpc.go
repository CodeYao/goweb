package carpc

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type Args struct {
	Args []interface{}
}
type Result struct {
	Value bool
}

type Address struct {
}

func (address *Address) JudgeAddress(args *Args, result *Result) error {
	fmt.Println("chenyao ******************* successfully", args.Args[0])
	result.Value = false
	return nil
}

func RegRPC() {
	var address = new(Address)
	rpc.Register(address)
	rpc.HandleHTTP() //将Rpc绑定到HTTP协议上。
	fmt.Println("启动服务...")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("服务已停止!")
}
