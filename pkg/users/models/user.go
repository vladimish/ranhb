package models

type User struct {
	Id                   int64
	Group                string
	LastAction           string
	LastActionValue      int
	IsPrivacyAccepted    bool
	IsTermsOfUseAccepted bool
	RegisterTime         int64
	LastActionTime       int64
	ActionsAmount        int
}
