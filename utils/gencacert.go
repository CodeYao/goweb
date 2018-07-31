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
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/tjfoc/gmsm/sm2"
)

func genCACert(privateKey string) []byte {
	switch Cfg.Cert.KeyType {
	case "sm2":
		fmt.Println("use sm2 cert")
		return genSM2CACert(privateKey)
	case "ecdsa":
		fmt.Println("use ECDSA cert")
		return genECDSACACert(privateKey)
	default:
		fmt.Println("err key type!")
	}
	return nil
}

func genSM2CACert(privateKey string) []byte {
	fmt.Println("==========SM2============")
	//ca读私钥
	privKey, err := sm2.ReadPrivateKeyFromMem([]byte(privateKey), nil)
	if err != nil {
		fmt.Println("read priv key err:", err)
	}

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
	ok, err := sm2.CreateCertificateToMem(&template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		fmt.Println("sm create cert err")
		return nil
	}

	fmt.Println("create node cert!")

	return ok
}

func genECDSACACert(data string) []byte {
	fmt.Println("==========ECDSA============")
	privKey, err := ecdsaPrivKeyFromMen([]byte(data))
	if err != nil {

		fmt.Println("read ecdsa priv err:", err)
		return nil
	}
	template := x509Template()
	template.IsCA = true
	_, certbyte, err := ecdsaCert(&template, &template, &privKey.PublicKey, privKey)
	//_, err := ecdsaCert("./ca/ca.pem", &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		fmt.Println("gen ca cert err:", err)
		return nil
	}
	fmt.Println("create ca cert success!")

	// req = ./req/priv?_req.pem
	//s := fmt.Sprintf("./req/%s_req.pem", strings.TrimSuffix(PrivK, ".pem"))

	return certbyte
}
