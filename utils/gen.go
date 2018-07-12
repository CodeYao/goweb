// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/tjfoc/gmsm/sm2"
)

var oidExtensionSubjectKeyId = []int{2, 5, 29, 14}
var PrivK string = "key.pem"

func genCert(path string) {
	switch Cfg.Cert.KeyType {
	case "sm2":
		genSM2Cert(path)
	case "ecdsa":
		genECDSACert(path)
	default:
		fmt.Println("err key type!")
	}
}

func genSM2Cert(path string) {
	fmt.Println("==========SM2============")
	//ca读私钥
	privKey, err := sm2.ReadPrivateKeyFromPem("./conf/ca/key.pem", nil)
	if err != nil {
		fmt.Println("read priv key err:", err)
	}
	//检查ca证书是否存在
	if !CheckIsExist("./conf/ca/ca.pem") {

		var ip []net.IP
		for _, v := range Cfg.Cert.IPAddress {
			ip = append(ip, net.ParseIP(v))

		}

		template := sm2.Certificate{
			SerialNumber: big.NewInt(-1),
			Subject: pkix.Name{
				CommonName:   Cfg.Cert.CommonName,
				Organization: Cfg.Cert.Organization,
				Country:      Cfg.Cert.Country,
				ExtraNames: []pkix.AttributeTypeAndValue{
					{
						Type:  []int{2, 5, 4, 42},
						Value: "Gopher",
					},

					{
						Type:  []int{2, 5, 4, 6},
						Value: "NL",
					},
				},
			},
			//			NotBefore: time.Unix(int64(Cfg.Cert.NotBefore), 0),
			NotBefore: time.Unix(time.Now().Unix()-int64(3600*24*Cfg.Cert.NotBefore), 0),
			NotAfter:  time.Unix(time.Now().Unix()+int64(3600*24*Cfg.Cert.NotAfter), 0),

			SignatureAlgorithm: sm2.SignatureAlgorithm(Cfg.Cert.SM2SignatureAlgorithm),

			SubjectKeyId: []byte{1, 2, 3, 4},
			KeyUsage:     sm2.KeyUsageCertSign,

			ExtKeyUsage:        []sm2.ExtKeyUsage{sm2.ExtKeyUsageClientAuth, sm2.ExtKeyUsageServerAuth},
			UnknownExtKeyUsage: []asn1.ObjectIdentifier{[]int{1, 2, 3}, []int{2, 59, 1}},

			BasicConstraintsValid: true,
			IsCA: true,

			OCSPServer:            []string{"http://ocsp.example.com"},
			IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

			DNSNames:       Cfg.Cert.DNSNames,
			EmailAddresses: Cfg.Cert.EmailAddresses,

			//		IPAddresses: []net.IP{net.ParseIP("2001:4860:0:2001::68")},
			IPAddresses: ip,

			PolicyIdentifiers:   []asn1.ObjectIdentifier{[]int{1, 2, 3}},
			PermittedDNSDomains: Cfg.Cert.PermittedDNSDomains,

			CRLDistributionPoints: Cfg.Cert.CRLDistributionPoints,

			ExtraExtensions: []pkix.Extension{
				{
					Id:    []int{1, 2, 3, 4},
					Value: []byte("extra extension"),
				},

				{
					Id:       oidExtensionSubjectKeyId,
					Critical: false,
					Value:    []byte{0x04, 0x04, 4, 3, 2, 1},
				},
			},
		}

		fmt.Println("create ca cert!")

		// _, err = sm2.WritePrivateKeytoPem("privv.pem", privKey, nil)
		// if err != nil {
		// 	fmt.Println("----------------------------")
		// }
		template.IsCA = true
		//生成ca证书
		ok, _ := sm2.CreateCertificateToPem("conf/ca/ca.pem", &template, &template, &privKey.PublicKey, privKey)
		if !ok {
			fmt.Println("sm create cert err")
			return
		}

	}

	fmt.Println("create node cert!")
	//privKey, err := sm2.ReadPrivateKeyFromPem(PrivK, nil)
	//if err != nil {
	//	fmt.Println("read priv key err:", err)
	//}
	s := fmt.Sprintf(path+"/%s_req.pem", strings.TrimSuffix(PrivK, ".pem"))

	//读取证书生成请求
	req, err := sm2.ReadCertificateRequestFromPem(s)
	if err != nil {
		fmt.Println("read req err:", err)
		return
	}
	//设置证书路径
	cert, err := sm2.ReadCertificateFromPem("conf/ca/ca.pem")
	tmp := sm2CsrToCert(req)
	//filepath := "./" + path + "/" + PrivK
	//s = fmt.Sprintf("%s_cert.pem", strings.TrimSuffix(filepath, ".pem"))
	s = "./" + path + "/cert.pem"
	pub := req.PublicKey

	v, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("err pub key type")
		return
	}

	smPub := &sm2.PublicKey{
		Curve: v.Curve,
		X:     v.X,
		Y:     v.Y,
	}

	//！！！！！！
	/**
	*s:路径
	*tmp:证书格式的节点请求
	*cert:ca证书
	*smPub:节点公钥
	*privKey:ca私钥
	 */
	ok, err = sm2.CreateCertificateToPem(s, tmp, cert, smPub, privKey)
	if !ok {
		fmt.Println("create cert err!", err)
		return
	}

	cert1, err := sm2.ReadCertificateFromPem(s)
	if err != nil {
		fmt.Println("read cert err:", err)
		return
	}
	err = cert.CheckSignature(cert1.SignatureAlgorithm, cert1.RawTBSCertificate, cert1.Signature)
	if err != nil {
		fmt.Println("check signature err:", err)
		return
	}

	fmt.Println("********success!********")
	fmt.Println("==========SM2============")
}

func genECDSACert(path string) {
	fmt.Println("==========ECDSA============")
	privKey, err := ecdsaPrivKeyFromPem("./ca/key.pem")
	if err != nil {

		fmt.Println("read ecdsa priv err:", err)
		return
	}
	if !CheckIsExist("./ca/ca.pem") {
		template := x509Template()
		template.IsCA = true
		_, err := ecdsaCert("./ca/ca.pem", &template, &template, &privKey.PublicKey, privKey)
		if err != nil {
			fmt.Println("gen ca cert err:", err)
			return
		}
		fmt.Println("create ca cert success!")
	}

	// req = ./req/priv?_req.pem
	s := fmt.Sprintf("./req/%s_req.pem", strings.TrimSuffix(PrivK, ".pem"))

	req, err := parseECDSAReq(s)
	if err != nil {
		fmt.Println("parse ecdsa req err:", err)
		return
	}

	//req -> x509.Certificate
	tmp := ecdsaCsrToCert(req)
	//s = priv?_cert.pem
	//filepath := "./" + path + "/" + PrivK
	//s = fmt.Sprintf("%s_cert.pem", strings.TrimSuffix(filepath, ".pem"))
	s = "./" + path + "/cert.pem"
	cert, err := parseECDSACert("./ca/ca.pem")
	//	fmt.Println("cert:", cert)
	if err != nil {
		fmt.Println("parse ecdsa cert err:", err)
		return
	}
	//gen node cert!
	pub := req.PublicKey
	//	fmt.Printf("pub:%x\n", pub)
	v, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("key type err")
	}
	cert1, err := ecdsaCert(s, tmp, cert, v, privKey)
	if err != nil {
		fmt.Println("create cert err!", err)
		return
	}
	//parse priv?_cert.pem
	// certb, err := parseCert(s)
	// if err != nil {
	// 	fmt.Println("read cert err:", err)
	// 	return
	// }

	//get node cert
	// cert1, err := x509.ParseCertificate(certb)
	// if err != nil {
	// 	fmt.Println("x509 parse cert err:", err)
	// 	return
	// }
	err = cert.CheckSignature(cert1.SignatureAlgorithm, cert1.RawTBSCertificate, cert1.Signature)
	if err != nil {
		fmt.Println("check signature err:", err)
		fmt.Println("==========ECDSA============")
		return
	}

	fmt.Println("create node cert success!")
	fmt.Println("==========ECDSA============")
}

func ecdsaCert(fileName string, tmp, parent *x509.Certificate, pub *ecdsa.PublicKey, priv interface{}) (*x509.Certificate, error) {
	b, err := x509.CreateCertificate(rand.Reader, tmp, parent, pub, priv)
	if err != nil {
		fmt.Println("create cert err:", err)
		return nil, err
	}

	certFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: b})

	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, err
	}

	return x509Cert, nil

}

func x509Template() x509.Certificate {
	var ip []net.IP
	for _, v := range Cfg.Cert.IPAddress {
		ip = append(ip, net.ParseIP(v))
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(-1),
		Subject: pkix.Name{
			CommonName:   Cfg.Cert.CommonName,
			Organization: Cfg.Cert.Organization,
			Country:      Cfg.Cert.Country,
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  []int{2, 5, 4, 42},
					Value: "Gopher",
				},
				// This should override the Country, above.
				{
					Type:  []int{2, 5, 4, 6},
					Value: "NL",
				},
			},
		},
		NotBefore:             time.Unix(time.Now().Unix()-int64(3600*24*Cfg.Cert.NotBefore), 0),
		NotAfter:              time.Unix(time.Now().Unix()+int64(3600*24*Cfg.Cert.NotAfter), 0),
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SignatureAlgorithm(Cfg.Cert.ECDSASignatureAlgorithm),
		DNSNames:              Cfg.Cert.DNSNames,
		EmailAddresses:        Cfg.Cert.EmailAddresses,
		IPAddresses:           ip,
		PermittedDNSDomains:   Cfg.Cert.PermittedDNSDomains,
		CRLDistributionPoints: Cfg.Cert.CRLDistributionPoints,

		SubjectKeyId: []byte{1, 2, 3, 4},

		OCSPServer:            []string{"http://ocsp.example.com"},
		IssuingCertificateURL: []string{"http://crt.example.com/ca1.crt"},

		PolicyIdentifiers: []asn1.ObjectIdentifier{[]int{1, 2, 3}},

		ExcludedDNSDomains: []string{"bar.example.com"},

		ExtraExtensions: []pkix.Extension{
			{
				Id:    []int{1, 2, 3, 4},
				Value: []byte("extra extension"),
			},
			// This extension should override the SubjectKeyId, above.
			{
				Id:       oidExtensionSubjectKeyId,
				Critical: false,
				Value:    []byte{0x04, 0x04, 4, 3, 2, 1},
			},
		},
	}
	return template
}

func sm2CsrToCert(req *sm2.CertificateRequest) *sm2.Certificate {

	if req == nil {
		fmt.Println("req = nil")
		return nil
	}
	var ip []net.IP
	for _, v := range Cfg.Cert.IPAddress {
		ip = append(ip, net.ParseIP(v))

	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	sm2 := &sm2.Certificate{
		Raw: req.Raw,
		//		RawTBSCertificate:       req.RawTBSCertificate,
		RawSubjectPublicKeyInfo: req.RawSubjectPublicKeyInfo,
		//		RawIssuer:               req.RawIssuer,
		Signature:          req.Signature,
		SignatureAlgorithm: req.SignatureAlgorithm,
		PublicKeyAlgorithm: req.PublicKeyAlgorithm,
		PublicKey:          req.PublicKey,
		Version:            req.Version,
		SerialNumber:       serialNumber,
		Subject:            req.Subject,
		NotBefore:          time.Unix(time.Now().Unix()-int64(3600*24*Cfg.Cert.NotBefore), 0),
		NotAfter:           time.Unix(time.Now().Unix()+int64(3600*24*Cfg.Cert.NotAfter), 0),

		DNSNames:              Cfg.Cert.DNSNames,
		EmailAddresses:        Cfg.Cert.EmailAddresses,
		IPAddresses:           ip,
		PermittedDNSDomains:   Cfg.Cert.PermittedDNSDomains,
		CRLDistributionPoints: Cfg.Cert.CRLDistributionPoints,
	}

	return sm2

}

func ecdsaCsrToCert(req *x509.CertificateRequest) *x509.Certificate {

	if req == nil {
		fmt.Println("req = nil")
		return nil
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	var ip []net.IP
	for _, v := range Cfg.Cert.IPAddress {
		ip = append(ip, net.ParseIP(v))

	}
	x509 := &x509.Certificate{
		Raw: req.Raw,
		//		RawTBSCertificate:       req.RawTBSCertificate,
		RawSubjectPublicKeyInfo: req.RawSubjectPublicKeyInfo,
		//		RawIssuer:               req.RawIssuer,
		Signature:          req.Signature,
		SignatureAlgorithm: req.SignatureAlgorithm,
		PublicKeyAlgorithm: req.PublicKeyAlgorithm,
		PublicKey:          req.PublicKey,
		Version:            req.Version,
		SerialNumber:       serialNumber,
		Subject:            req.Subject,
		NotBefore:          time.Unix(time.Now().Unix()-int64(3600*24*Cfg.Cert.NotBefore), 0),
		NotAfter:           time.Unix(time.Now().Unix()+int64(3600*24*Cfg.Cert.NotAfter), 0),

		DNSNames:              req.DNSNames,
		EmailAddresses:        req.EmailAddresses,
		IPAddresses:           ip,
		PermittedDNSDomains:   Cfg.Cert.PermittedDNSDomains,
		CRLDistributionPoints: Cfg.Cert.CRLDistributionPoints,
		//	PermittedDNSDomains: req.PermittedDNSDomains,
		//		CRLDistributionPoints: .CRLDistributionPoints,
	}

	return x509

}
