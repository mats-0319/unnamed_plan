package main

import (
	utilsdb "github.com/mats0319/unnamed_plan/server/internal/utils/init_db"
)

func main() {
	dbConfig := utilsdb.DefaultConfig()
	db := utilsdb.InitDB(dbConfig)

	db.Create(defaultUser)
	//db.Create(testUser)
	//db.Create(testNote)
}
