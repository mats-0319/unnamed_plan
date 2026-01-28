package model

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Note struct {
	mdb_model.Model
	NoteID      string `gorm:"unique;not null"` // 使用其他字段计算获得，可用于保证新增接口幂等性
	WriterID    uint   // user id
	WriterName  string // user nickname(at that time)
	IsAnonymous bool   // 是否匿名，仅前端展示使用
	Title       string
	Content     string `gorm:"not null"`
}

func NewNote(writerID uint, writerName string, isAnonymous bool, title string, content string) *Note {
	noteIns := &Note{
		WriterID:    writerID,
		WriterName:  writerName,
		IsAnonymous: isAnonymous,
		Title:       title,
		Content:     content,
	}

	noteBytes := fmt.Sprintf(`"writer id":%d,"writer name":%s,"is anonymous":%t,"title":%s,"content":%s`,
		noteIns.WriterID, noteIns.WriterName, noteIns.IsAnonymous, noteIns.Title, noteIns.Content)

	noteIns.NoteID = utils.HmacSHA256[string](noteBytes) // 保证新增接口幂等性

	return noteIns
}
