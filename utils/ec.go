package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func parseCert(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	p, _ := pem.Decode(b)
	return p.Bytes, nil

}
func ecdsaPrivKeyFromMen(data []byte) (*ecdsa.PrivateKey, error) {

	p, _ := pem.Decode(data)

	key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	v, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("parse pkcs8 priv key err")
	}
	return v, nil

}
func ecdsaPrivKeyFromPem(fileName string) (*ecdsa.PrivateKey, error) {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	p, _ := pem.Decode(b)

	key, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	v, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("parse pkcs8 priv key err")
	}
	return v, nil

}

func parseECDSAReq(file []byte) (*x509.CertificateRequest, error) {

	//	fmt.Println("file name:", fileName)
	// file, err := os.Open(fileName)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer file.Close()
	// b, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	p, err := pem.Decode(file)
	if err != nil {
		fmt.Println(err)
	}

	//	fmt.Printf("rrrrrr:%x\n", p.Bytes)
	return x509.ParseCertificateRequest(p.Bytes)

}

func parseECDSACert(fileName string) (*x509.Certificate, error) {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	p, _ := pem.Decode(b)
	if err != nil {
		fmt.Println(err)
	}

	return x509.ParseCertificate(p.Bytes)
}

var (
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}

	oidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
)

type privateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
}

func marshalEcdsaUnecryptedPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	var r pkcs8
	var priv privateKey
	var algo pkix.AlgorithmIdentifier

	algo.Algorithm = oidPublicKeyECDSA
	algo.Parameters.Class = 0
	algo.Parameters.Tag = 6
	algo.Parameters.IsCompound = false
	algo.Parameters.FullBytes, _ = asn1.Marshal(oidNamedCurveP256)

	priv.Version = 1
	priv.NamedCurveOID = oidNamedCurveP256
	priv.PublicKey = asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)}
	priv.PrivateKey = key.D.Bytes()
	r.Version = 0
	r.Algo = algo
	r.PrivateKey, _ = asn1.Marshal(priv)
	return asn1.Marshal(r)
}
