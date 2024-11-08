package models

import (
	"time"

)


type User struct {
	Id int64
	Username string
	FirstName string
	LastName string
	Birthday time.Time
	Email string
	Password string
	PasswordHashed []byte
	Role string
}

type Token struct {
	Id int
	AccessToken string
}