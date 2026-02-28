package mdb

import (
	"encoding/json"
	"os"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type config struct {
	DSN          string `json:"dsn"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}

func Initialize() {
	configIns, err := getDBConfig()
	if err != nil {
		os.Exit(1)
	}

	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	}
	if utils.IsTestMode {
		gormConfig.NamingStrategy = schema.NamingStrategy{TablePrefix: "t_"}
	}

	db, err := gorm.Open(postgres.Open(configIns.DSN), gormConfig)
	if err != nil {
		mlog.Error("open db failed", mlog.Field("error", err))
		os.Exit(1)
	}

	// set conns
	{
		sqlDB, err := db.DB()
		if err != nil {
			mlog.Error("get sql db failed", mlog.Field("error", err))
			os.Exit(1)
		}

		sqlDB.SetMaxIdleConns(configIns.MaxIdleConns)
		sqlDB.SetMaxOpenConns(configIns.MaxOpenConns)
	}

	dal.SetDefault(db)

	logStr := "> DB init."
	if utils.IsTestMode {
		logStr += " [test mode]"
	}
	mlog.Info(logStr)
}

func getDBConfig() (*config, error) {
	bytes := mconfig.GetConfigItem(utils.ConfigID_DB)

	conf := &config{}
	err := json.Unmarshal(bytes, conf)
	if err != nil {
		mlog.Error("deserialize db config failed", mlog.Field("error", err))
		return nil, err
	}

	return conf, nil
}
