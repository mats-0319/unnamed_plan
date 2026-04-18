package middleware

import (
	"fmt"
	"testing"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func TestAccessToken(t *testing.T) {
	token := GenerateApiAccessToken("test user name")

	ctx := &mhttp.Context{AccessToken: token}
	if err := VerifyAccessToken(ctx); err != nil {
		t.Error(err.String())
	}
}

func TestAccessTokenExceptions(t *testing.T) {
	// invalid structure
	tokenWithoutPoint := "aStrWithoutPoint"
	ctx := &mhttp.Context{AccessToken: tokenWithoutPoint}

	if err := verifyAccessToken(ctx); err.Error() != utils.ErrInvalidAccessToken().Error() {
		t.Error("should fail")
	}

	// tamper hash
	tokenTamper := GenerateApiAccessToken("test user name")
	tokenTamper = tokenTamper[:len(tokenTamper)-1] + "g"

	ctx.AccessToken = tokenTamper
	if err := verifyAccessToken(ctx); err.Error() != utils.ErrWrongAccessTokenHash().Error() {
		t.Error("should fail")
	}

	// wrong token type
	tokenWrongType := generateToken(&Token{
		UserName:   "test user name",
		Type:       TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Hour).UnixMilli(),
	})

	ctx.AccessToken = tokenWrongType
	if err := verifyAccessToken(ctx); err.Error() != utils.ErrInvalidAccessToken().Error() {
		t.Error("should fail")
	}

	// token expired
	tokenExpired := generateToken(&Token{
		UserName:   "test user name",
		Type:       TokenType_ApiAccessToken,
		ExpireTime: time.Now().Add(-1 * time.Minute).UnixMilli(),
	})

	ctx.AccessToken = tokenExpired
	if err := verifyAccessToken(ctx); err.Error() != utils.ErrAccessTokenExpired().Error() {
		t.Error("should fail")
	}
}

func TestMfaToken(t *testing.T) {
	// test clear expired token, including basic function test
	lastMinute := time.Now().Add(-1 * time.Minute).UnixMilli()
	for i := range 1000 {
		t := &Token{
			UserName:   fmt.Sprintf("test user name %d", i),
			Type:       TokenType_MFAToken,
			ExpireTime: lastMinute,
		}

		token := generateToken(t)

		NewMFAToken(t.UserName, token, lastMinute)
	}

	// 测试基本功能
	validToken := GenerateMFAToken("test user name")
	tokenCountBeforeClear := len(mtm.Data)

	if err := VerifyMFAToken("test user name", validToken); err != nil {
		t.Error(err.String())
	}

	tokenCountAfterClear := 0
	for range 3 {
		tokenCountAfterClear = len(mtm.Data)
		t.Logf("time: %v, token count: %d", time.Now(), tokenCountAfterClear)

		time.Sleep(time.Microsecond * 10)
	}

	if tokenCountBeforeClear <= 1000 || tokenCountAfterClear >= 1000 {
		t.Error(fmt.Sprintf("token clear unexpected, token count from %d to %d",
			tokenCountBeforeClear, tokenCountAfterClear))
	}
}
