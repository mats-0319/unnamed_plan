package config

import (
	"encoding/json"
	"log"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
)

type Config struct {
	HMACKey               string `json:"hmac_key"`
	AccessTokenExpireHour int    `json:"access_token_expire_hour"`
	MFATokenExpireMinute  int    `json:"mfa_token_expire_minute"`
	TOTPKeyExpireMinute   int    `json:"totp_key_expire_minute"`
}

var conf = &Config{}

func Init() {
	configBytes := mconfig.GetConfigItem("3e6fe66d-32bb-46b7-9597-8de23a969706")

	if err := json.Unmarshal(configBytes, conf); err != nil {
		log.Fatalln("parse config failed, error: ", err)
	}
}

func GetServerConfig() *Config {
	return conf
}
