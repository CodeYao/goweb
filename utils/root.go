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
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"wutongMG/goweb/models"

	"github.com/tjfoc/gmsm/sm2"
)

//var cfgFile string = "./cmd/conf.yaml"
var Cfg ToolConf

func init() {
	//Cfg.Cert.KeyType = "sm2"
	//Cfg.Cert.CommonName = "test.example.com"
	//Cfg.Cert.Organization = []string{"TEST", "TEST1"}
	//Cfg.Cert.Country = []string{"China"}
	Cfg.Cert.NotBefore = 100
	//Cfg.Cert.NotAfter = 1000
	//Cfg.Cert.DNSNames = []string{"10.1.3.150", "10.1.3.150"}
	Cfg.Cert.PermittedDNSDomains = []string{".example.com", "example.com"}
	Cfg.Cert.CRLDistributionPoints = []string{"http://crl1.example.com/ca1.cr1", "http://cr12.example.com/ca1.crl"}
	Cfg.Cert.SM2SignatureAlgorithm = 18
	Cfg.Cert.ECDSASignatureAlgorithm = 10
	// viper.SetConfigFile(cfgFile)
	// viper.AutomaticEnv()
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println(err)
	// }

	// if err := viper.Unmarshal(&Cfg); err != nil {
	// 	fmt.Println(err)
	// }
}

func Creatdir(dir string) {
	//dir := "./gzFiles2"
	exist := CheckIsExist(dir)

	if exist {
		//fmt.Printf("has dir![%v]\n", dir)
	} else {
		//fmt.Printf("no dir![%v]\n", dir)
		// 创建文件夹
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			//fmt.Printf("mkdir success!\n")
		}
	}
}

func CheckPrivAndCert(certdata []byte, privPubdata []byte) bool {
	cert, err := sm2.ReadCertificateFromMem(certdata)
	if err != nil {
		fmt.Println(err)
		return false
	}
	privPub, err := sm2.ReadPublicKeyFromMem(privPubdata, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	certpub := cert.PublicKey.(*ecdsa.PublicKey)
	smPub := &sm2.PublicKey{
		Curve: certpub.Curve,
		X:     certpub.X,
		Y:     certpub.Y,
	}
	fmt.Println(privPub.X)
	fmt.Println(smPub.X)
	fmt.Println(privPub.Y)
	fmt.Println(smPub.Y)

	if privPub.X.Cmp(smPub.X) == 0 && privPub.Y.Cmp(smPub.Y) == 0 {
		fmt.Println("证书和私钥成功匹配")
		return true
	}
	return false
}

func GetEcdsaPubKey(reqdata string) []byte {
	req, err := parseECDSAReq([]byte(reqdata))
	if err != nil {
		panic(err)
	}
	pubkey := req.PublicKey.(*ecdsa.PublicKey)
	pubB, err := x509.MarshalPKIXPublicKey(pubkey)
	pub := &pem.Block{
		Type:  "ECDSA PUBLIC KEY",
		Bytes: pubB,
	}
	// s := "/pubkey.pem"
	// //设置公钥路径
	// s = "conf/" + ipAddress + s
	ok := pem.EncodeToMemory(pub)

	return ok
}

func VerifySM2(cert string) error {
	_, err := sm2.ReadCertificateFromMem([]byte(cert))
	return err
}

func GetSm2PubKeyFromPrivKey(privKey string) ([]byte, error) {
	priv, err := sm2.ReadPrivateKeyFromMem([]byte(privKey), nil)
	if err != nil {
		return nil, err
	}
	pub := priv.PublicKey
	ok, err := sm2.WritePublicKeytoMem(&pub, nil)
	if err != nil {
		fmt.Println(err)
	}
	return ok, nil
}

func GetSm2PubKey(reqdata string) []byte {
	req, err := sm2.ReadCertificateRequestFromMem([]byte(reqdata))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pubkey := req.PublicKey.(*ecdsa.PublicKey)
	smPub := &sm2.PublicKey{
		Curve: pubkey.Curve,
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	// s := "/pubkey.pem"
	// //设置公钥路径
	// s = "conf/" + ipAddress + s
	ok, err := sm2.WritePublicKeytoMem(smPub, nil)
	if err != nil {
		fmt.Println(err)
	}
	return ok
}

func GenerateCACert(notAfter int, ipAddress []string, country []string, organization []string, commonName string, privKey string, KeyType string) []byte {
	Cfg.Cert.KeyType = KeyType
	Cfg.Cert.NotAfter = notAfter
	Cfg.Cert.DNSNames = ipAddress
	if len(organization) > 0 {
		Cfg.Cert.Organization = organization
	}
	if len(country) > 0 {
		Cfg.Cert.Country = country
	}
	if commonName != "" {
		Cfg.Cert.CommonName = commonName
	}
	//Creatdir("conf/ca")

	//fmt.Println(Cfg.Cert)
	return genCACert(privKey)
}

func GenerateCert(notAfter int, ipAddress []string, country []string, organization []string, commonName string, reqData string, KeyType string) []byte {
	Cfg.Cert.KeyType = KeyType
	//设定证书的keytype与CA的相同
	caInfo, err := models.QueryData("ca where enabled = 'enabled'")
	if err != nil {
		models.Errorf("Get CA info error: %v", err)
		return nil
	}
	Cfg.Cert.KeyType = caInfo[0]["keytype"]
	Cfg.Cert.NotAfter = notAfter
	Cfg.Cert.DNSNames = ipAddress
	if len(organization) > 0 {
		Cfg.Cert.Organization = organization
	}
	if len(country) > 0 {
		Cfg.Cert.Country = country
	}
	if commonName != "" {
		Cfg.Cert.CommonName = commonName
	}
	//Creatdir("conf/ca")
	//Creatdir("conf/req")
	// if !CheckIsExist("conf/ca/key.pem") {
	// 	genKey("conf/ca")
	// }
	//reqpath := strings.Join(ipAddress, ",")
	//path := "conf/" + reqpath
	//Creatdir(path)
	//genKey(path)

	//genCetReq(path)
	models.Debugf("%v", Cfg.Cert)
	return genCert(reqData)
	//os.RemoveAll("conf/req")
}
func GetCert(certbyte []byte) sm2.Certificate {
	cert, err := sm2.ReadCertificateFromMem(certbyte)
	if err != nil {
		panic(err)
	}
	return *cert
}

//path:想要吊销的证书的列表
func RevokedCertificates(path []string) {
	caInfo, err := models.QueryData("ca where enabled = 'enabled'")
	if err != nil {
		panic(err)
	}
	cert, err := sm2.ReadCertificateFromMem([]byte(caInfo[0]["cacert"]))
	// fmt.Println("***chenyao***", caInfo)
	if err != nil {
		panic(err)
	}
	privKey, err := sm2.ReadPrivateKeyFromMem([]byte(caInfo[0]["caprivkey"]), nil)
	if err != nil {
		fmt.Println("read priv key err:", err)
	}
	var revokedCerts []pkix.RevokedCertificate
	for _, v := range path {
		peer_cert, err := sm2.ReadCertificateFromMem([]byte(v))
		if err != nil {
			fmt.Println(err)
		}
		peer_revoke := pkix.RevokedCertificate{
			SerialNumber:   peer_cert.SerialNumber,
			RevocationTime: time.Now(),
			Extensions:     peer_cert.Extensions,
		}
		revokedCerts = append(revokedCerts, peer_revoke)
	}

	now := time.Now()
	dd, _ := time.ParseDuration("24000h")
	expir := now.Add(dd)
	crl, err := cert.CreateCRL(rand.Reader, privKey, revokedCerts, now, expir)
	if err != nil {
		panic(err)
	}
	block := &pem.Block{
		Type:  "CRL",
		Bytes: crl,
	}
	file, err := os.Create("./conf/crl.pem")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = pem.Encode(file, block)
	if err != nil {
		panic(err)
	}

	rcrl, err := readClr("./conf/crl.pem")
	if err != nil {
		panic(err)
	}
	err = cert.CheckCRLSignature(rcrl)
	if err != nil {
		panic(err)
	}
	//fmt.Println("吊销成功")
}

// func EasyGen(ipAddress []string) {
// 	Cfg.Cert.DNSNames = ipAddress
// 	Creatdir("conf/ca")
// 	if !CheckIsExist("conf/ca/key.pem") {
// 		genKey("conf/ca")
// 	}
// 	reqpath := strings.Join(ipAddress, ",")
// 	path := "conf/" + reqpath
// 	Creatdir(path)
// 	genKey(path)
// 	genCetReq(path)
// 	genCert(path)
// }

func OneTouch(KeyType string) []string {
	Cfg.Cert.KeyType = KeyType
	var pub []byte
	key := genKey()
	req := genCetReq(key)
	if KeyType == "sm2" {
		pub = GetSm2PubKey(string(req))
	} else if KeyType == "ecdsa" {
		pub = GetEcdsaPubKey(string(req))
	}

	return []string{string(key), string(pub), string(req)}
}

func ReadECDSACertFromMen(certdata []byte) (*x509.Certificate, error) {
	p, err := pem.Decode(certdata)
	if err != nil {
		fmt.Println(err)
	}

	return x509.ParseCertificate(p.Bytes)
}

func GetCAKey(KeyType string) (string, string) {
	Cfg.Cert.KeyType = KeyType
	var pub []byte
	key := genKey()
	req := genCetReq(key)
	if KeyType == "sm2" {
		pub = GetSm2PubKey(string(req))
	} else if KeyType == "ecdsa" {
		pub = GetEcdsaPubKey(string(req))
	}

	return string(key), string(pub)
}

func readClr(path string) (*pkix.CertificateList, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode CRL")
	}
	return sm2.ParseCRL(block.Bytes)

}
