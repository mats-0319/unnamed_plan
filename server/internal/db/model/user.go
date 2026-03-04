package model

type User struct {
	Model
	UserName string `gorm:"unique;not null"` // login name, can't modify
	Nickname string `gorm:"unique;not null"` // display name
	Password string `gorm:"not null"`
	IsAdmin  bool

	Enable2FA bool
	TotpKey   string `gorm:"size:16"`

	// last login, use 'user.UpdateAt'
}
