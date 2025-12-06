package model

import "gorm.io/gorm"

type Model struct {
	ID        uint           `gorm:"primaryKey"` // gorm tag不区分大小写
	CreatedAt int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
