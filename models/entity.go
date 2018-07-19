package models

import "encoding/xml"

type Account struct {
	AccountId    string `json:accountId`
	Password     string `json:password`
	Organization string `json:organization`
	AccountLevel string `json:accountLevel`
	Enable       string `json:enable`
}

type Resource struct {
	XMLName    xml.Name `xml:"resource"`
	Dbhostsip  string   `xml:"dbhostsip"`
	Dbusername string   `xml:"dbusername"`
	Dbpassowrd string   `xml:"dbpassword"`
	Dbname     string   `xml:"dbname"`
}

type CertVO struct {
	//IpReq     string `json:ipReq`
	CertName     string `json:certName`
	CertDay      string `json:certDay`
	IpAdderss    string `json:ipAddress`
	Country      string `json:country`
	Organization string `json:organization`
	CommonName   string `json:commonName`
	State        string `json:state`
}
