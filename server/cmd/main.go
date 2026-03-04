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

	brm := backup.NewBRManager(&backup.UserBR{}, &backup.NoteBR{})

	wg := &sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	go autoBackup(ctx, wg, brm)
	go waitSignal(ctx, wg, brm) // wait SIGUSR1 / SIGUSR2

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
	h.AddHandler(uriPrefix+api.URI_ListUser, handlers.ListUser, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_ModifyUser, handlers.ModifyUser, middleware.VerifyAccessToken)

	// note
	h.AddHandler(uriPrefix+api.URI_CreateNote, handlers.CreateNote, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_ListNote, handlers.ListNote)
	h.AddHandler(uriPrefix+api.URI_ModifyNote, handlers.ModifyNote, middleware.VerifyAccessToken)
	h.AddHandler(uriPrefix+api.URI_DeleteNote, handlers.DeleteNote, middleware.VerifyAccessToken)

	return h
}

func autoBackup(ctx context.Context, wg *sync.WaitGroup, brm *backup.BRManager) {
	time.Sleep(500 * time.Millisecond)

	mlog.Info("> Goroutine: Auto Backup Start.")
	defer wg.Done()

	for {
		// timestamp, unit: second
		now := time.Now().Unix()
		tomorrowZero := (now/86400 + 1) * 86400 // 0时区，次日0点
		remain := tomorrowZero - now

		select {
		case <-ctx.Done():
			mlog.Info("> Goroutine: Auto Backup Exit.")
			return
		case <-time.After(time.Duration(remain) * time.Second):
			mlog.Info("> Auto Backup Start ...")
			err := brm.Backup()
			if err != nil {
				mlog.Error("auto backup failed", mlog.Field("error", err))
			} else {
				mlog.Info("> Auto Backup Done.")
			}
		}
	}
}

func waitSignal(ctx context.Context, wg *sync.WaitGroup, brm *backup.BRManager) {
	time.Sleep(500 * time.Millisecond)

	mlog.Info("> Goroutine: Wait Signal Start.")
	defer wg.Done()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)

	for {
		select {
		case <-ctx.Done():
			mlog.Info("> Goroutine: Wait Signal Exit.")
			return
		case sig := <-ch:
			mlog.Info("> Received Signal: " + sig.String())

			switch sig {
			case syscall.SIGUSR1:
				err := brm.Backup()
				if err != nil {
					mlog.Error("backup failed", mlog.Field("error", err))
				} else {
					mlog.Info("> Backup Done.")
				}
			case syscall.SIGUSR2:
				err := brm.Recover()
				if err != nil {
					mlog.Error("recover failed", mlog.Field("error", err))
				} else {
					mlog.Info("> Recover Done.")
				}
			}
		}
	}
}
