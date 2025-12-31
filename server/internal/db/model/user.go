package model

type User struct {
	Model
	UserName  string `gorm:"unique;not null"`  // login name
	Nickname  string `gorm:"unique;not null"`  // display name
	Password  string `gorm:"size:64;not null"` // hex(sha256(hex(sha256(text)), salt))
	Salt      string `gorm:"size:10;not null"`
	TotpKey   string `gorm:"size:16"` // 允许为空，需要设置后启用
	IsAdmin   bool
	LastLogin int64 // timestamp, unit: milli
}
