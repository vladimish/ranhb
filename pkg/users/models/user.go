package models

type User struct {
	Id                int64
	Group             string
	LastAction        string
	LastActionValue   int
	IsPrivacyAccepted bool
}
