package middleware

import (
	"testing"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func TestAccessToken(t *testing.T) {
	token := GenerateApiAccessToken("user name 1")
	t.Log("> Token: ", token)

	ctx := &mhttp.Context{AccessToken: token}
	if err := VerifyAccessToken(ctx); err != nil {
		t.Error(err.String())
	}

	t.Logf("> Verified, user id: %s", ctx.UserName)
}
