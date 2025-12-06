package model

type User struct {
	Model
	Name       string `gorm:"unique;not null"`  // login name
	Nickname   string `gorm:"not null"`         // display name
	Password   string `gorm:"size:64;not null"` // sha256(sha256(text), salt)
	Salt       string `gorm:"size:10;not null"`
	TotpKey    string `gorm:"size:16"` // 允许为空，需要设置后启动
	IsLocked   bool
	Permission uint8  // 权限等级
	Creator    string // 创建人，user.name
	LastLogin  int64  // timestamp, unit: milli
}
