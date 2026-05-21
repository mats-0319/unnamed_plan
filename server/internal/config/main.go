package mconfig

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Level    string          `json:"level"`
	Internal *InternalConfig `json:"internal"`
	Items    []*ConfigItem   `json:"items"`
}

type InternalConfig struct {
	HTTPServerPort int    `json:"http_server_port"`
	DBDSN          string `json:"db_dsn"` // data resource name
}

type ConfigItem struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Json json.RawMessage `json:"json"`
}

var conf = &Config{}

func Initialize(configInitFunc ...func()) {
	confBytes, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("> Read config file failed, path: 'config.json', error: ", err)
	}

	if err := json.Unmarshal(confBytes, conf); err != nil {
		log.Fatalln("> Json unmarshal failed, error: ", err)
	}

	for _, f := range configInitFunc {
		f()
	}
}

func GetLevel() string {
	return conf.Level
}

func GetInternalConfig() *InternalConfig {
	return conf.Internal
}

func GetConfigItem(id string) json.RawMessage {
	var res json.RawMessage
	for _, v := range conf.Items {
		if v.ID == id {
			res = v.Json
			break
		}
	}

	return res
}
