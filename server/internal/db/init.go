package mdb

import (
	"encoding/json"
	"os"

	mconf "github.com/mats0319/unnamed_plan/server/internal/config"
	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	DSN          string `json:"dsn"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}

func init() {
	conf, err := getDBConfig()
	if err != nil {
		mlog.Log("get db config failed", mlog.Field("error", err))
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		mlog.Log("open db failed", mlog.Field("error", err))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		mlog.Log("get sql db failed", mlog.Field("error", err))
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	dal.SetDefault(db)

	mlog.Log("> DB init.")
}

func getDBConfig() (*config, error) {
	bytes := mconf.GetConfigItem(mconst.UID_DB)

	conf := &config{}
	if err := json.Unmarshal(bytes, conf); err != nil {
		mlog.Log("deserialize db config failed", mlog.Field("error", err))
		return nil, err
	}

	return conf, nil
}
