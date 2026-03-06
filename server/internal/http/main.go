package mhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type config struct {
	Port string `json:"port"`
}

// StartServer is blocked
func StartServer(handler *Handler) {
	configIns, err := getHttpConfig()
	if err != nil {
		os.Exit(1)
	}

	handler.displayRegisteredUri()

	// 因为允许手机访问需要设置web ip为`192.168.xxx`，即本机内网ipv4地址，
	// 这里为了支持`127.0.0.1`、内网ip等多种格式，使用`0.0.0.0`
	addr := fmt.Sprintf("0.0.0.0:%s", configIns.Port)
	mlog.Info("> Listening at: " + addr)

	err = http.ListenAndServe(addr, handler)
	if err != nil {
		mlog.Error("handlers listen and serve failed", mlog.Field("error", err))
	}
}

func getHttpConfig() (*config, error) {
	jsonBytes := mconfig.GetConfigItem(utils.ConfigID_Http)

	res := &config{}
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		mlog.Error("deserialize gateway config failed", mlog.Field("error", err))
		return nil, err
	}

	return res, nil
}
