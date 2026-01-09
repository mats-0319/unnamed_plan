package model

import "fmt"

type Note struct {
	Model
	NoteID      string `gorm:"unique;not null"` // 展示用
	WriterID    uint   // user id
	WriterName  string // user nickname(at that time)
	IsAnonymous bool   // 是否匿名，仅前端展示使用
	Title       string
	Content     string `gorm:"not null"`
}

func (n *Note) Serialize() string {
	return fmt.Sprintf(`"writer id":%d,"writer name":%s,"is anonymous":%t,"title":%s,"content":%s`,
		n.WriterID, n.WriterName, n.IsAnonymous, n.Title, n.Content)
}
