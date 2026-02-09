package middleware

import (
	"testing"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func TestAccessToken(t *testing.T) {
	token := GenAccessToken(1001)
	t.Log("> Token: ", token)

	ctx := &mhttp.Context{AccessToken: token}
	err := VerifyAccessToken(ctx)
	if err != nil {
		t.Error(err.String())
	}

	t.Logf("> Verified, user id: %d", ctx.User)
}
