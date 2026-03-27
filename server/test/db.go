package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	utilsdb "github.com/mats0319/unnamed_plan/server/internal/utils/init_db"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
	"gorm.io/gorm"
)

var db *gorm.DB

func createTable() {
	dropTable()

	if err := db.Migrator().CreateTable(model.ModelList...); err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	db.Create(presetUser())
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

func presetUser() []*model.User {
	pwdFromWeb := utils.CalcSHA256("123456")

	return []*model.User{
		{
			UserName:  "admin",
			Nickname:  "admin",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			IsAdmin:   true,
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user",
			Nickname:  "user",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user_with_totp",
			Nickname:  "user_with_totp",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			Enable2FA: true,
			TotpKey:   "NVQXE2LP", // base32 of 'mario'
		},
	}
}
