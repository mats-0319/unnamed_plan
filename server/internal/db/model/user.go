package model

type User struct {
	Model
	Name      string `gorm:"unique;not null"`  // login name
	Nickname  string `gorm:"unique;not null"`  // display name
	Password  string `gorm:"size:64;not null"` // hex(sha256(sha256(text), salt))
	Salt      string `gorm:"size:10;not null"`
	TotpKey   string // 允许为空，需要设置后启动
	IsAdmin   bool
	LastLogin int64 // timestamp, unit: milli
}
