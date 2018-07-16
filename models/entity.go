package models

type Account struct {
	AccountId    string `json:accountId`
	Password     string `json:password`
	Organization string `json:organization`
	AccountLevel string `json:accountLevel`
	Enable       string `json:enable`
}
