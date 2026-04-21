package mdb

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfig struct {
	DSN          string
	IsTestMode   bool
	MaxIdleConns int
	MaxOpenConns int
}

func NewConfig(dsn string, isTestMode bool, maxIdleConns int, maxOpenConns int) *DBConfig {
	return &DBConfig{
		DSN:          dsn,
		IsTestMode:   isTestMode,
		MaxIdleConns: maxIdleConns,
		MaxOpenConns: maxOpenConns,
	}
}

func DefaultConfig(isTestMode bool) *DBConfig {
	return NewConfig("host=115.190.167.134 user=mario password=123456 dbname=test_cloud port=5432 sslmode=disable",
		isTestMode, 10, 100)
}

func InitDB(config *DBConfig) *gorm.DB {
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Error),
		SkipDefaultTransaction: true,
	}
	if config.IsTestMode {
		gormConfig.NamingStrategy = schema.NamingStrategy{TablePrefix: "t_"}
	}

	db, err := gorm.Open(postgres.Open(config.DSN), gormConfig)
	if err != nil {
		log.Fatalln("open db failed, error: ", err)
	}

	// set conns
	{
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalln("get sql db failed, error: ", err)
		}

		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}

	dal.SetDefault(db)

	return db
}
