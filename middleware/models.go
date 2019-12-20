package middleware

import "time"

type Model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Permission struct {
	Model
	Name  string  `json:"name"`
	Roles []*Role `gorm:"many2many:role_permissions;"`
}

type Role struct {
	Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
	Accounts    []Account     `gorm:"foreignkey:RoleID"`
}

type Account struct {
	Model
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
