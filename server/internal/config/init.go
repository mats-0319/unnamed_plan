package mconfig

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Level string        `json:"level"`
	Items []*ConfigItem `json:"items"`
}

type ConfigItem struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Json json.RawMessage `json:"json"`
}

var conf = &Config{}

func init() {
	confBytes, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("> Read config file failed, path: 'config.json', error: ", err)
	}

	if err = json.Unmarshal(confBytes, conf); err != nil {
		log.Fatalln("> Json unmarshal failed, error: ", err)
	}
}

func GetLevel() string {
	return conf.Level
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
