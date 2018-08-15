package ca_grpcserver

import (
	pb "ca/goweb/ca_grpc"
	"ca/goweb/models"
	"ca/goweb/utils"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/tjfoc/gmsm/sm2"
	"google.golang.org/grpc"
)

const (
	port = ":9092"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
// func (s *server) JudgeAddress(ctx context.Context, in *pb.AddressRequest) (*pb.AddressReply, error) {
// 	return &pb.AddressReply{Message: "Hello " + in.Addr}, nil
// }

func (s *server) VerifyAddress(ctx context.Context, in *pb.AddressRequest) (*pb.IsPermissionReply, error) {
	addresslist := models.QueryData("address where enabled = 'enabled'")
	for _, v := range addresslist {
		if in.Addr == v["address"] {
			//fmt.Println("chenyao**************true")
			fmt.Println("***token address verify success***")
			return &pb.IsPermissionReply{IsPermission: true}, nil
		}
	}
	fmt.Println("***token address verify failed***")
	//fmt.Println("chenyao**************false")
	return &pb.IsPermissionReply{IsPermission: false}, nil
}

func (s *server) GetAddressList(ctx context.Context, in *pb.Empty) (*pb.AddressList, error) {
	addresstable := models.QueryData("codeaddress where enabled = 'enabled'")
	addresslist := []string{}
	for _, v := range addresstable {
		addresslist = append(addresslist, v["address"])
	}
	fmt.Println("***get code address success***")
	return &pb.AddressList{Addresslist: addresslist}, nil
}

func (s *server) VerifyCert(ctx context.Context, in *pb.Cert) (*pb.IsPermissionReply, error) {
	cacertinfo := models.QueryData("ca where enabled = 'enabled'")
	cacertstr := cacertinfo[0]["cacert"]
	in.Keytype = strings.ToLower(in.Keytype)
	if in.Keytype == "sm2" {
		cacert, err := sm2.ReadCertificateFromMem([]byte(cacertstr))
		if err != nil {
			fmt.Println("parse sm2 cert err:", err)
			panic(err)
		}
		var peercert sm2.Certificate
		err = json.Unmarshal(in.Cert, &peercert)
		//peercert, err := sm2.ReadCertificateFromMem(in.Cert)
		if err != nil {
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		err = cacert.CheckSignature(peercert.SignatureAlgorithm, peercert.RawTBSCertificate, peercert.Signature)
		if err != nil {
			fmt.Println("check signature err:", err)
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		fmt.Println("check cert success!")
		fmt.Println("==========sm2============")
		return &pb.IsPermissionReply{IsPermission: true}, nil
	} else if in.Keytype == "ecdsa" {
		cacert, err := utils.ReadECDSACertFromMen([]byte(cacertstr))
		if err != nil {
			fmt.Println("parse ecdsa cert err:", err)
			panic(err)
		}
		var peercert x509.Certificate
		err = json.Unmarshal(in.Cert, &peercert)
		//peercert, err := utils.ReadECDSACertFromMen(in.Cert)
		if err != nil {
			fmt.Println("parse ecdsa cert err:", err)
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		err = cacert.CheckSignature(peercert.SignatureAlgorithm, peercert.RawTBSCertificate, peercert.Signature)
		if err != nil {
			fmt.Println("check signature err:", err)
			fmt.Println("==========ECDSA============")
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		fmt.Println("check cert success!")
		fmt.Println("==========ECDSA============")
		return &pb.IsPermissionReply{IsPermission: true}, nil
	}
	return &pb.IsPermissionReply{IsPermission: false}, nil

}

func CAGrpcRun() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}
