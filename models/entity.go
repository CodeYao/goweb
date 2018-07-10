package models

type Account struct {
	AccountId       string `json:accountId`
	AccountPassword string `json:accountPassword`
	AccountLevel    string `json:accountLevel`
}
