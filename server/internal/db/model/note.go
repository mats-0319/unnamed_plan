package model

type Note struct {
	Model
	NoteID      string `gorm:"unique;not null"` // 展示用
	WriterID    uint   // user id
	WriterName  string // user nickname
	IsAnonymous bool   // 是否匿名，仅前端展示使用
	Title       string
	Content     string `gorm:"not null"`
}
