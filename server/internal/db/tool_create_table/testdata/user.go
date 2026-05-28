package testdata

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

var TestUsers = []*model.User{
	NewUser("admin", "", true, true, "5SSFNNEJUENPCCKP"),
	NewUser("user", "", false, false, ""),
}

func NewUser(userName string, nickname string, isAdmin bool, enable2FA bool, totpKey string) *model.User {
	pwdSHA256 := utils.CalcSHA256("123456")
	pwdArgon2 := password.GeneratePassword(pwdSHA256)

	if len(nickname) < 1 {
		nickname = userName
	}

	return &model.User{
		UserName:  userName,
		Nickname:  nickname,
		Password:  pwdArgon2,
		IsAdmin:   isAdmin,
		EnableMFA: enable2FA,
		TOTPKey:   totpKey,
	}
}
