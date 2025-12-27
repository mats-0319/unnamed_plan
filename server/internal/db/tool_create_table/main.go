package main

import (
	"fmt"
	"log"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "host=115.190.167.134 user=mario password=123456 dbname=test_cloud port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // print all sql
	})
	if err != nil {
		log.Fatalln("open db failed, error: ", err)
	}

	// gorm新增记录，如果一组记录中的某一个已经存在，则后续记录不会处理，所以每次更新测试数据时我们会删除表并重建
	// 删除表: 如果在dbeaver中已经打开一个表查看其数据，则重建后原本的表无法查看，重新打开即可
	err = db.Migrator().DropTable(model.ModelList...)
	if err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}

	err = db.Migrator().CreateTable(model.ModelList...)
	if err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	tableNames, err := db.Migrator().GetTables()
	if err != nil {
		log.Fatalln("get table names failed, error: ", err)
	}

	// 修改sequence，设置初始值
	for _, v := range tableNames {
		sequence := fmt.Sprintf("%s_id_seq", v)
		db.Exec(fmt.Sprintf("alter sequence %s restart with 1001;", sequence))
	}

	db.Create(defaultUsers)
	db.Create(testNotes)
}
