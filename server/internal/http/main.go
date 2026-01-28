package mhttp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type config struct {
	Port string `json:"port"`
}

// StartServer is blocked
func StartServer(handler *Handler) {
	handler.config = getConfig()

	handler.supportedUri()

	addr := fmt.Sprintf("%s:%s", "0.0.0.0", handler.config.Port)
	mlog.Log("> Listening at: " + addr)

	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatalln("handlers listen and serve failed", err)
	}
}

func getConfig() *config {
	jsonBytes := mconfig.GetConfigItem("22c4db2e-06d3-4d6a-b43f-c42aa6f57d15")

	res := &config{}
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		log.Fatalln("get gateway config failed", err)
	}

	return res
}
