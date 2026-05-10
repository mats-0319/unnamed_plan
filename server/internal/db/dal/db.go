package dal

import "gorm.io/gorm"

func DB() *gorm.DB {
	return Q.db
}
