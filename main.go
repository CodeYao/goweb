package main

import (
	"wutongMG/goweb/ca_grpcserver"
	"wutongMG/goweb/controllers"
)

func main() {
	go ca_grpcserver.CAGrpcRun()

	controllers.RunWeb()
}
