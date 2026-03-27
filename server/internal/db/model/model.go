package model

import (
	"github.com/google/uuid"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"` // gorm tag不区分大小写
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	BackupAt  int64          // 不适合使用自动编辑
}

func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		newID, err := uuid.NewV7()
		if err != nil {
			mlog.Error("generate uuid failed", mlog.Field("error", err))
			return err
		}

		m.ID = newID
	}

	return nil
}

// ModelList use for gorm tools
var ModelList = []any{
	&User{},
	&Note{},
	&FlipGameScore{},
}
