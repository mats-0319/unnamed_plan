package middleware

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func TestAccessToken(t *testing.T) {
	token := GenerateAPIAccessToken("test user name")

	ctx := &mhttp.Context{AccessToken: token}
	if err := VerifyAccessToken(ctx); err != nil {
		t.Error(err.String())
	}
}

func TestAccessTokenExceptions(t *testing.T) {
	// invalid structure
	tokenWithoutPoint := "aStrWithoutPoint"
	ctx := &mhttp.Context{AccessToken: tokenWithoutPoint}

	if verifyAccessToken(ctx).Code != utils.ErrInvalidAccessToken().Code {
		t.Error("position 1")
	}

	// tamper hash
	tokenTamper := GenerateAPIAccessToken("test user name")
	tokenTamper = tokenTamper[:len(tokenTamper)-1] + "g"

	ctx.AccessToken = tokenTamper
	if verifyAccessToken(ctx).Code != utils.ErrWrongAccessTokenHash().Code {
		t.Error("position 2")
	}

	// wrong token type
	tokenWrongType := generateToken(&Token{
		UserName:   "test user name",
		Type:       TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Hour).UnixMilli(),
	})

	ctx.AccessToken = tokenWrongType
	if verifyAccessToken(ctx).Code != utils.ErrInvalidAccessToken().Code {
		t.Error("position 3")
	}

	// token expired
	tokenExpired := generateToken(&Token{
		UserName:   "test user name",
		Type:       TokenType_APIAccessToken,
		ExpireTime: time.Now().Add(-1 * time.Minute).UnixMilli(),
	})

	ctx.AccessToken = tokenExpired
	if verifyAccessToken(ctx).Code != utils.ErrAccessTokenExpired().Code {
		t.Error("position 4")
	}
}

func TestMFAToken(t *testing.T) {
	// test clear expired token, including basic function test
	lastMinute := time.Now().Add(-1 * time.Minute).UnixMilli()
	for i := range 1000 {
		tokenIns := &Token{
			UserName:   fmt.Sprintf("test user name %d", i),
			Type:       TokenType_MFAToken,
			ExpireTime: lastMinute,
		}

		token := generateToken(tokenIns)

		NewMFAToken(tokenIns.UserName, token, lastMinute)
	}

	// 测试基本功能
	validToken := GenerateMFAToken("test user name")
	tokenCountBeforeClear := len(mtm.Data)

	if _, err := VerifyMFAToken(validToken); err != nil {
		t.Error(err.String())
	}

	tokenCountAfterClear := len(mtm.Data)
	for i := 0; i < 3; {
		newLength := len(mtm.Data)
		if newLength != tokenCountAfterClear {
			i++
			tokenCountAfterClear = newLength
			t.Logf("time: %v, token count: %d", time.Now(), tokenCountAfterClear)
		} else if newLength == 0 {
			break
		}
	}

	if tokenCountBeforeClear <= 1000 || tokenCountAfterClear >= 1000 {
		t.Error(fmt.Sprintf("token clear unexpected, token count from %d to %d",
			tokenCountBeforeClear, tokenCountAfterClear))
	}
}

// 测试同时只运行一个实例，如果已经有一个实例在运行，忽略启动新实例的行为
func TestRunOneInstanceAtOnce(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	var wg sync.WaitGroup
	var lock atomic.Bool

	f := func() {
		log.Println("> New Function Call")

		if lock.CompareAndSwap(false, true) {
			log.Println("| Start New Goroutine")

			wg.Go(func() {
				defer lock.Store(false)

				time.Sleep(time.Millisecond * 500)
				log.Println("| Exit New Goroutine")
			})
		}
	}

	for range 10 {
		f()
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
}
