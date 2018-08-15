package main

import (
	"encoding/json"
	"fmt"
	"国密算法/gmsm/sm2"
)

func main() {
	cert, err := sm2.ReadCertificateFromPem("C:/Users/chenyao/Desktop/cert.pem")
	if err != nil {
		panic(err)
	}
	pageJson, _ := json.Marshal(cert)
	fmt.Println(string(pageJson))
}
