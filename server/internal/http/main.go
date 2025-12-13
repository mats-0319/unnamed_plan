package mhttp

import (
	"encoding/json"
	"log"
	"net/http"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type config struct {
	Port string `json:"port"`
}

// StartServer is blocked
func StartServer(handler *Handler) {
	configIns := getConfig()

	handler.supportedUri()

	addr := "127.0.0.1:" + configIns.Port
	mlog.Log("> Listening at: " + addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalln("handlers listen and serve failed", err)
	}
}

func getConfig() *config {
	jsonBytes := mconfig.GetConfigItem(mconst.UID_Http)

	res := &config{}
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		log.Fatalln("get gateway config failed", err)
	}

	return res
}
