package middleware

import "time"

type model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Permission struct {
	model
	Name  string `json:"name"`
	Roles []Role
}

type Role struct {
	model
	Name        string `json:"name"`
	Description string `json:"description"`
	Accounts    []Account
	Permissions []Permission
}

type Account struct {
	model
	UserName  string `json:"username" sql:"index"`
	Email     string `json:"email" sql:"index"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Enable    bool   `json:"enable"`
	Locked    bool   `json:"locked"`
	Role      Role
	RoleID    int64
}
