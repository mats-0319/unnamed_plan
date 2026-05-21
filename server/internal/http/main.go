package mhttp

import (
	"fmt"
	"log/slog"
	"net/http"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

// StartServer is blocked
func StartServer(handler *Handler) {
	handler.DisplayRegisteredURI()

	// 因为允许手机访问需要设置web ip为`192.168.xxx`(即本机内网ipv4地址)，
	// 这里为了支持`127.0.0.1`、内网ip等多种格式，使用`0.0.0.0`
	addr := fmt.Sprintf("0.0.0.0:%d", mconfig.GetInternalConfig().HTTPServerPort)
	mlog.Info("> Listening at: " + addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		mlog.Error("handlers listen and serve failed", slog.Any("error", err))
	}
}
