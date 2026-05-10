package testdata

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func PresetUser() []*model.User {
	pwdSHA256 := utils.CalcSHA256("123456")

	return []*model.User{
		{
			UserName:  "admin",
			Nickname:  "admin",
			Password:  password.GeneratePassword(pwdSHA256),
			IsAdmin:   true,
			EnableMFA: false,
			TOTPKey:   "",
		},
		{
			UserName:  "user",
			Nickname:  "user",
			Password:  password.GeneratePassword(pwdSHA256),
			EnableMFA: false,
			TOTPKey:   "",
		},
		{
			UserName:  "user_with_totp",
			Nickname:  "user_with_totp",
			Password:  password.GeneratePassword(pwdSHA256),
			EnableMFA: true,
			TOTPKey:   "NVQXE2LP", // base32 of 'mario'
		},
	}
}
