package config

import (
	"encoding/json"
	"log"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
)

type Config struct {
	Address string `json:"address"`
}

var ConfigIns *Config

func init() {
	ConfigIns = getConfig()
}

func getConfig() *Config {
	jsonBytes := mconfig.GetConfigItem(mconst.UID_UserServer)

	res := &Config{}
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		log.Fatalln("get gateway config failed", err)
	}

	return res
}
