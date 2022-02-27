package models

type UserAccount struct {
	ID           string
	Username     string
	PasswordHash string
	Email        *string
}
