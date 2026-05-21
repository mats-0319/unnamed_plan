package mfa

import (
	"encoding/base32"
	"time"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func GenerateTOTPKey(userName string, expireMinute int) string {
	keyBytes := utils.GenerateRandomBytes[[]byte](10)
	keyBase32 := base32.StdEncoding.EncodeToString(keyBytes)

	expireTime := time.Now().Add(time.Duration(expireMinute) * time.Minute).UnixMilli()

	newTOTPKey(userName, keyBase32, expireTime)

	return keyBase32
}

func VerifyTOTPKey(userName string, code string) (key string, e *utils.Error) {
	return maintainTOTPKey(userName, code)
}
