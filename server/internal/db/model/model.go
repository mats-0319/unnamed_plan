package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuidv7()"` // gorm tag不区分大小写
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	BackupAt  int64          // 不适合使用自动编辑
}

// ModelList use for gorm tools
var ModelList = []any{
	&User{},
	&Note{},
	&FlipGameScore{},
}
