package main

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	utilsdb "github.com/mats0319/unnamed_plan/server/internal/utils/init_db"
)

func main() {
	dbConfig := utilsdb.DefaultConfig()
	db := utilsdb.InitDB(dbConfig)

	if err := db.Migrator().DropTable(model.ModelList...); err != nil {
		fmt.Println("drop table failed, err: ", err)
		return
	}

	if err := db.Migrator().CreateTable(model.ModelList...); err != nil && err.Error() != "insufficient arguments" {
		fmt.Println("create table failed, err: ", err)
		return
	}

	db.Create(defaultUser)
	db.Create(testFlipGameScore)
	//db.Create(testUser)
	//db.Create(testNote)
}
