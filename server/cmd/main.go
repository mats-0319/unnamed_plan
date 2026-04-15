package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers"
	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/backup"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func main() {
	mconfig.Initialize()
	mlog.Initialize()
	defer mlog.Close()
	mdb.Initialize()

	ctx, cancel := context.WithCancel(context.Background())
	brm := backup.NewBRManager(&backup.UserBR{}, &backup.NoteBR{}, &backup.FlipGameScore{})

	var wg sync.WaitGroup

	wg.Go(func() { autoBackup(ctx, brm) })
	wg.Go(func() { waitSignal(ctx, brm) }) // wait SIGUSR1 / SIGUSR2

	mhttp.StartServer(newHandler()) // blocked
	cancel()

	wg.Wait()
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	uriPrefix := "/api" // even use domain name like 'api.xxx.com/login', nginx will forward req

	// user
	h.AddHandler(uriPrefix+api.URI_Register, handlers.Register)
	h.AddHandler(uriPrefix+api.URI_Login, handlers.Login)
	h.AddHandler(uriPrefix+api.URI_LoginTotp, handlers.LoginTotp)
	h.AddHandler(uriPrefix+api.URI_ListUser, handlers.ListUser, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_ModifyUser, handlers.ModifyUser, middleware.VerifyAccessToken)

	// note
	h.AddHandler(uriPrefix+api.URI_CreateNote, handlers.CreateNote, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_ListNote, handlers.ListNote, middleware.OptionalVerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_ModifyNote, handlers.ModifyNote, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_DeleteNote, handlers.DeleteNote, middleware.VerifyAccessToken)

	// game score
	h.AddHandler(uriPrefix+api.URI_ListGameScore, handlers.ListGameScore)
	h.AddHandler(uriPrefix+api.URI_UploadGameScore, handlers.UploadGameScore, middleware.OptionalVerifyAccessToken)

	return h
}

func autoBackup(ctx context.Context, brm *backup.BRManager) {
	time.Sleep(500 * time.Millisecond)

	mlog.Info("> Goroutine: Auto Backup Start.")

	for {
		// timestamp, unit: second
		now := time.Now().Unix()
		tomorrowZero := (now/86400 + 1) * 86400 // 0时区，次日0点，即北京时间每天早上8点
		remain := tomorrowZero - now

		select {
		case <-ctx.Done():
			mlog.Info("> Goroutine: Auto Backup Exit.")
			return
		case <-time.After(time.Duration(remain) * time.Second):
			mlog.Info("> Auto Backup Start ...")
			brm.Backup()
			mlog.Info("> Auto Backup Done.")
		}
	}
}

func waitSignal(ctx context.Context, brm *backup.BRManager) {
	time.Sleep(500 * time.Millisecond)

	mlog.Info("> Goroutine: Wait Signal Start.")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		select {
		case <-ctx.Done():
			mlog.Info("> Goroutine: Wait Signal Exit.")
			return
		case sig := <-sigCh:
			mlog.Info("> Received Signal: " + sig.String())

			switch sig {
			case syscall.SIGUSR1:
				brm.Backup()
				mlog.Info("> Backup Done.")
			case syscall.SIGUSR2:
				brm.Recover()
				mlog.Info("> Recover Done.")
			}
		}
	}
}
