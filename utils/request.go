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
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"

	"github.com/tjfoc/gmsm/sm2"
)

var DNames []string
var EMailAddr []string
var IPAddr []net.IP
var PrivateKey string = "key.pem"

func genCetReq(key []byte) []byte {
	switch Cfg.Cert.KeyType {
	case "sm2":
		return genSM2Req(key)
	case "ecdsa":
		return genECDSAReq(key)
	default:
		fmt.Println("err key type!")
	}
	return nil
}

func genSM2Req(key []byte) []byte {
	privKey, err := sm2.ReadPrivateKeyFromMem(key, nil)
	if err != nil {
		fmt.Println(err)
	}

	DNames = Cfg.Cert.DNSNames
	EMailAddr = Cfg.Cert.EmailAddresses
	//	IPAddr = Cfg.Cert.IPAddress
	tmpReq := sm2.CertificateRequest{
		SignatureAlgorithm: sm2.SignatureAlgorithm(Cfg.Cert.SM2SignatureAlgorithm),
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Test"},
		},

		//		NotBefore: time.Unix(Cfg.Cert.NotBefore, 0),

		//		NotAfter:  time.Unix(Cfg.Cert.NotAfter, 0),
		//TODO:
		PublicKey:      privKey.PublicKey,
		DNSNames:       DNames,
		EmailAddresses: EMailAddr,
		IPAddresses:    IPAddr,
	}
	//s := fmt.Sprintf("./"+path+"/%s_req.pem", strings.TrimSuffix(PrivateKey, ".pem"))
	ok, err := sm2.CreateCertificateRequestToMem(&tmpReq, privKey)
	if err != nil {
		fmt.Println(err)
	}
	return ok
}

func genECDSAReq(data []byte) []byte {
	DNames = Cfg.Cert.DNSNames
	EMailAddr = Cfg.Cert.EmailAddresses

	priv, err := ecdsaPrivKeyFromMen(data)
	if err != nil {
		fmt.Println(err)
	}

	tmpReq := x509.CertificateRequest{
		SignatureAlgorithm: x509.SignatureAlgorithm(Cfg.Cert.ECDSASignatureAlgorithm),
		Subject: pkix.Name{
			CommonName:   Cfg.Cert.CommonName,
			Organization: Cfg.Cert.Organization,
		},

		PublicKey: priv.PublicKey,
		//TODO:
		DNSNames:       DNames,
		EmailAddresses: EMailAddr,
		IPAddresses:    IPAddr,
	}

	//	fmt.Printf("priv:%x\n", priv.PublicKey)
	der, err := x509.CreateCertificateRequest(rand.Reader, &tmpReq, priv)

	block := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: der,
	}
	return pem.EncodeToMemory(block)
}
