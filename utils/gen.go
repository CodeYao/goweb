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
	"time"
	"wutongMG/goweb/models"

	"github.com/tjfoc/gmsm/sm2"
)

var oidExtensionSubjectKeyId = []int{2, 5, 29, 14}
var PrivK string = "key.pem"

func genCert(reqdata string) []byte {
	switch Cfg.Cert.KeyType {
	case "sm2":
		models.Infof("use sm2 cert")
		return genSM2Cert(reqdata)
	case "ecdsa":
		models.Infof("use ECDSA cert")
		return genECDSACert(reqdata)
	default:
		models.Errorf("err key type!")
	}
	return nil
}

func genSM2Cert(reqData string) []byte {
	models.Infof("==========SM2============")
	//ca读私钥
	caInfo, err := models.QueryData("ca where enabled = 'enabled'")
	if err != nil {
		models.Errorf("get CA info error: %v", err)
	}
	privKey, err := sm2.ReadPrivateKeyFromMem([]byte(caInfo[0]["caprivkey"]), nil)
	//fmt.Println("*****chenyao*****", caInfo)
	if err != nil {
		models.Errorf("read priv key err: %v", err)
	}

	models.Infof("create node cert!")
	//privKey, err := sm2.ReadPrivateKeyFromPem(PrivK, nil)
	//if err != nil {
	//	fmt.Println("read priv key err:", err)
	//}
	//s := fmt.Sprintf(path+"/%s_req.pem", strings.TrimSuffix(PrivK, ".pem"))

	//读取证书生成请求
	req, err := sm2.ReadCertificateRequestFromMem([]byte(reqData))
	if err != nil {
		models.Errorf("read req err: %v", err)
		return nil
	}
	//设置证书路径
	cert, err := sm2.ReadCertificateFromMem([]byte(caInfo[0]["cacert"]))
	tmp := sm2CsrToCert(req)
	//filepath := "./" + path + "/" + PrivK
	//s = fmt.Sprintf("%s_cert.pem", strings.TrimSuffix(filepath, ".pem"))
	//s = "./" + path + "/cert.pem"
	tmp.Subject.CommonName = Cfg.Cert.CommonName
	tmp.Subject.Country = Cfg.Cert.Country
	tmp.Subject.Organization = Cfg.Cert.Organization
	pub := req.PublicKey

	v, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		models.Errorf("err pub key type")
		return nil
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
	newcert, err := sm2.CreateCertificateToMem(tmp, cert, smPub, privKey)
	if err != nil {
		models.Errorf("create cert err : %v", err)
		return nil
	}

	cert1, err := sm2.ReadCertificateFromMem(newcert)
	if err != nil {
		models.Errorf("read cert err: %v", err)
		return nil
	}
	err = cert.CheckSignature(cert1.SignatureAlgorithm, cert1.RawTBSCertificate, cert1.Signature)
	if err != nil {
		models.Errorf("check signature err: %v", err)
		return nil
	}

	models.Infof("********success!********")
	models.Infof("==========SM2============")
	return newcert
}

func genECDSACert(data string) []byte {
	models.Infof("==========ECDSA============")
	caInfo, err := models.QueryData("ca where enabled = 'enabled'")
	if err != nil {
		models.Errorf("get CA info err:%v", err)
		return nil
	}
	privKey, err := ecdsaPrivKeyFromMen([]byte(caInfo[0]["caprivkey"]))
	if err != nil {

		models.Errorf("read ecdsa priv err:%v", err)
		return nil
	}
	// req = ./req/priv?_req.pem
	//s := fmt.Sprintf("./req/%s_req.pem", strings.TrimSuffix(PrivK, ".pem"))

	req, err := parseECDSAReq([]byte(data))
	if err != nil {
		models.Fatalf("%v", err)
	}

	//req -> x509.Certificate
	tmp := ecdsaCsrToCert(req)
	//s = priv?_cert.pem
	//filepath := "./" + path + "/" + PrivK
	//s = fmt.Sprintf("%s_cert.pem", strings.TrimSuffix(filepath, ".pem"))
	//s = "./" + path + "/cert.pem"
	cert, err := parseECDSACert([]byte(caInfo[0]["cacert"]))
	//	fmt.Println("cert:", cert)
	if err != nil {
		models.Errorf("parse ecdsa cert err: %v", err)
		return nil
	}
	//gen node cert!
	pub := req.PublicKey
	//	fmt.Printf("pub:%x\n", pub)
	v, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		models.Errorf("key type err")
	}
	cert1, certbyte, err := ecdsaCert(tmp, cert, v, privKey)
	if err != nil {
		models.Errorf("create cert err : %v", err)
		return nil
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
		models.Errorf("check signature err: %v", err)
		models.Errorf("==========ECDSA============")
		return nil
	}

	models.Infof("create node cert success!")
	models.Infof("==========ECDSA============")
	return certbyte
}

func ecdsaCert(tmp, parent *x509.Certificate, pub *ecdsa.PublicKey, priv interface{}) (*x509.Certificate, []byte, error) {
	b, err := x509.CreateCertificate(rand.Reader, tmp, parent, pub, priv)
	if err != nil {
		fmt.Println("create cert err:", err)
		return nil, nil, err
	}

	certbyte := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})

	if err != nil {
		return nil, nil, err
	}

	x509Cert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, nil, err
	}

	return x509Cert, certbyte, nil

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
