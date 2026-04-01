package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	utilsdb "github.com/mats0319/unnamed_plan/server/internal/utils/init_db"
	"github.com/mats0319/unnamed_plan/server/test/testdata"
	"gorm.io/gorm"
)

var db *gorm.DB

func createTable() {
	dropTable()

	if err := db.Migrator().CreateTable(model.ModelList...); err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	db.Create(testdata.PresetUser())
}

func dropTable() {
	if db == nil {
		dbConfig := utilsdb.DefaultConfig()
		dbConfig.IsTestMode = true
		db = utilsdb.InitDB(dbConfig)
	}

	if err := db.Migrator().DropTable(model.ModelList...); err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}
}
