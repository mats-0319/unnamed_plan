package mfa

import (
	"time"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/token"
)

func GenerateMFAToken(userName string, expireMinute int) string {
	return token.SerializeToken(&token.Token{
		UserName:   userName,
		Type:       token.TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Duration(expireMinute) * time.Minute).UnixMilli(),
	})
}

func VerifyMFAToken(tokenStr string) (t *token.Token, e *utils.Error) {
	return token.DeserializeToken(tokenStr, token.TokenType_MFAToken)
}
