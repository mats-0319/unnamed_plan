package mdb

import (
	"encoding/json"
	"os"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
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
	configIns, err := getDBConfig()
	if err != nil {
		mlog.Log("get db config failed", mlog.Field("error", err))
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(configIns.DSN), &gorm.Config{
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

	sqlDB.SetMaxIdleConns(configIns.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configIns.MaxOpenConns)

	dal.SetDefault(db)

	mlog.Log("> DB init.")
}

func getDBConfig() (*config, error) {
	bytes := mconfig.GetConfigItem(mconst.UID_DB)

	conf := &config{}
	err := json.Unmarshal(bytes, conf)
	if err != nil {
		mlog.Log("deserialize db config failed", mlog.Field("error", err))
		return nil, err
	}

	return conf, nil
}
