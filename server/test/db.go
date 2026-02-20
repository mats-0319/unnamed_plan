package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/cmd/model"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func initDB() {
	dsn := "host=115.190.167.134 user=mario password=123456 dbname=test_cloud port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),    // print no sql
		NamingStrategy: schema.NamingStrategy{TablePrefix: "t_"}, // test table
	})
	if err != nil {
		log.Fatalln("open db failed, error: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("get sql db failed, error: ", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	dal.SetDefault(db)
}

func createTable() {
	dropTable()

	err := db.Migrator().CreateTable(model.ModelList...)
	if err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	db.Create(presetUser())
}

func dropTable() {
	if db == nil {
		initDB()
	}

	err := db.Migrator().DropTable(model.ModelList...)
	if err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}
}

func presetUser() []*model.User {
	pwdFromWeb := utils.CalcSHA256("123456")

	return []*model.User{
		{
			UserName:  "admin",
			Nickname:  "admin",
			Password:  utils.GeneratePwdHash(pwdFromWeb),
			IsAdmin:   true,
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user",
			Nickname:  "user",
			Password:  utils.GeneratePwdHash(pwdFromWeb),
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user_with_totp",
			Nickname:  "user_with_totp",
			Password:  utils.GeneratePwdHash(pwdFromWeb),
			Enable2FA: true,
			TotpKey:   "NVQXE2LP", // base32 of 'mario'
		},
	}
}
