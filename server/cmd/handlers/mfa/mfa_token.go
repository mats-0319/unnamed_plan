package mfa

import (
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/token"
)

func GenerateMFAToken(userName string, expireMinute int) string {
	tokenIns := &token.Token{
		UserName:   userName,
		Type:       token.TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Duration(expireMinute) * time.Minute).UnixMilli(),
	}

	newMFAToken(userName, tokenIns.ExpireTime)

	return token.SerializeToken(tokenIns)
}

func VerifyMFAToken(token string) (t *token.Token, e *utils.Error) {
	t, e = maintainMFAToken(token)
	if e != nil {
		mlog.Error(e.String())
		return
	}

	return
}
