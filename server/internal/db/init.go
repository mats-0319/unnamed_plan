package mdb

import (
	"encoding/json"
	"log/slog"
	"os"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/flag"
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

	dbConfig := NewConfig(configIns.DSN, flag.IsTestMode, configIns.MaxIdleConns, configIns.MaxOpenConns)
	_ = InitDB(dbConfig)

	logStr := "> DB init."
	if flag.IsTestMode {
		logStr += " [test mode]"
	}
	mlog.Info(logStr)
}

func getDBConfig() (*config, error) {
	bytes := mconfig.GetConfigItem(utils.ConfigID_DB)

	conf := &config{}
	if err := json.Unmarshal(bytes, conf); err != nil {
		mlog.Error("deserialize db config failed", slog.Any("error", err))
		return nil, err
	}

	return conf, nil
}
