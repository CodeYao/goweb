package ca_grpcserver

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"net"
	"strings"
	pb "wutongMG/goweb/ca_grpc"
	"wutongMG/goweb/controllers"
	"wutongMG/goweb/models"
	"wutongMG/goweb/utils"

	"github.com/tjfoc/gmsm/sm2"
	"google.golang.org/grpc"
)

var port string = controllers.WebConfig.GrpcPort

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
// func (s *server) JudgeAddress(ctx context.Context, in *pb.AddressRequest) (*pb.AddressReply, error) {
// 	return &pb.AddressReply{Message: "Hello " + in.Addr}, nil
// }

func (s *server) VerifyAddress(ctx context.Context, in *pb.AddressRequest) (*pb.IsPermissionReply, error) {
	addresslist, err := models.QueryData("address where enabled = 'enabled'")
	if err != nil {
		models.Errorf("get address list error: %s", err)
		return nil, err
	}
	for _, v := range addresslist {
		if in.Addr == v["address"] {
			//fmt.Println("chenyao**************true")
			models.Infof("token address verify success")
			//fmt.Println("***token address verify success***")
			return &pb.IsPermissionReply{IsPermission: true}, nil
		}
	}
	models.Infof("token address verify failed")
	//fmt.Println("***token address verify failed***")
	//fmt.Println("chenyao**************false")
	return &pb.IsPermissionReply{IsPermission: false}, nil
}

func (s *server) GetAddressList(ctx context.Context, in *pb.Empty) (*pb.AddressList, error) {
	addresstable, err := models.QueryData("codeaddress where enabled = 'enabled'")
	if err != nil {
		models.Errorf("get codeaddress table error: %s", err)
		return nil, err
	}
	addresslist := []string{}
	for _, v := range addresstable {
		addresslist = append(addresslist, v["address"])
	}
	models.Infof("get code address success")
	return &pb.AddressList{Addresslist: addresslist}, nil
}

func (s *server) VerifyCert(ctx context.Context, in *pb.Cert) (*pb.IsPermissionReply, error) {
	cacertinfo, err := models.QueryData("ca where enabled = 'enabled'")
	if err != nil {
		models.Errorf("get CA info error: %s", err)
		return nil, err
	}
	cacertstr := cacertinfo[0]["cacert"]
	in.Keytype = strings.ToLower(in.Keytype)
	if in.Keytype == "sm2" {
		cacert, err := sm2.ReadCertificateFromMem([]byte(cacertstr))
		if err != nil {
			models.Errorf("parse sm2 cert err: %s", err)
			//panic(err)
		}
		var peercert sm2.Certificate
		err = json.Unmarshal(in.Cert, &peercert)
		//peercert, err := sm2.ReadCertificateFromMem(in.Cert)
		if err != nil {
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		err = cacert.CheckSignature(peercert.SignatureAlgorithm, peercert.RawTBSCertificate, peercert.Signature)
		if err != nil {
			models.Errorf("check signature err: %s", err)
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		models.Infof("check cert success!")
		models.Infof("==========sm2============")
		return &pb.IsPermissionReply{IsPermission: true}, nil
	} else if in.Keytype == "ecdsa" {
		cacert, err := utils.ReadECDSACertFromMen([]byte(cacertstr))
		if err != nil {
			models.Errorf("parse ecdsa cert err: %s", err)
			panic(err)
		}
		var peercert x509.Certificate
		err = json.Unmarshal(in.Cert, &peercert)
		//peercert, err := utils.ReadECDSACertFromMen(in.Cert)
		if err != nil {
			models.Errorf("parse ecdsa cert err:", err)
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		err = cacert.CheckSignature(peercert.SignatureAlgorithm, peercert.RawTBSCertificate, peercert.Signature)
		if err != nil {
			models.Errorf("check signature err: %s", err)
			models.Errorf("==========ECDSA============")
			return &pb.IsPermissionReply{IsPermission: false}, err
		}
		models.Infof("check cert success!")
		models.Infof("==========ECDSA============")
		return &pb.IsPermissionReply{IsPermission: true}, nil
	}
	return &pb.IsPermissionReply{IsPermission: false}, nil

}

func CAGrpcRun() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		models.Fatalf("grpc failed to listen [%s]", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)
}
