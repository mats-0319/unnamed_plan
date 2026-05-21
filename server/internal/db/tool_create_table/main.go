package main

import (
	"fmt"

	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
)

func main() {
	db := mdb.InitDB(mdb.DefaultDSN, 10, 100)

	if err := db.Migrator().DropTable(model.ModelList...); err != nil {
		fmt.Println("drop db table failed, err: ", err)
		return
	}

	if err := db.Migrator().CreateTable(model.ModelList...); err != nil {
		fmt.Println("create db table failed, err: ", err)
		return
	}

	db.Create(defaultUser)

	db.Create(testUser)
	db.Create(testNote)
	db.Create(testFlipGameScore)
}
