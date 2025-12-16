package main

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=115.190.167.134 user=mario password=123456 dbname=test_cloud port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("open db failed, error: ", err)
	}

	models := []interface{}{
		&model.User{},
	}

	// gorm新增记录，如果一组记录中的某一个已经存在，则后续记录不会处理，所以每次更新测试数据时我们会删除表并重建
	// 删除表: 如果在dbeaver中已经打开一个表查看其数据，则重建后原本的表无法查看，重新打开即可
	if err = db.Migrator().DropTable(models...); err != nil {
		log.Fatalln("drop db table failed, error: ", err)
	}

	if err = db.Migrator().CreateTable(models...); err != nil {
		log.Fatalln("create db table failed, error: ", err)
	}

	db.Create(defaultUsers)
}

var defaultUsers = []*model.User{
	{
		UserName: "mats0319",
		Nickname: "Mario",
		Password: "15549caf9adafbe5543a26d51bcadba9029bbe5c871e8f195627a1a5fd77673c",
		Salt:     "9gIAEwTUge",
		TotpKey:  "5SSFNNEJUENPCCKP",
		IsAdmin:  true,
	}, {
		UserName: "admin",
		Nickname: "admin",
		Password: "a797a2163a2186de1d45e633b2e3a58bbb0eb3eb323fac22034fa27c067685cf",
		Salt:     "HGxHMnVXJS",
		TotpKey:  "",
		IsAdmin:  true,
	}, {
		UserName: "user",
		Nickname: "user",
		Password: "68ecd7703f2cb36105a51e42947d8cd000e101123b1f07133fe593557ebf1b22",
		Salt:     "3rWYu5vMNI",
		TotpKey:  "",
	},
}
