package models

// User model represents a registered user.
type User struct {
	ID       uint   `json:"id"   gorm:"not null; primaryKey; autoIncrement; comment:user id"`
	Name     string `json:"name" gorm:"not null; index; comment:user name"`
	Password string `json:"-"    gorm:"not null; comment:user password"`
}
