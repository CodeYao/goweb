package main

import (
	"ca/goweb/ca_grpcserver"
	"ca/goweb/controllers"
)

func main() {
	go ca_grpcserver.CAGrpcRun()

	controllers.RunWeb()
}
