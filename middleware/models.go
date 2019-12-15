package middleware

import "time"

type model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

/*
JWT claims struct
*/
type Token struct {
	model
	AccountID uint64 `json:"accountID"`
}

type Account struct {
	model
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}
