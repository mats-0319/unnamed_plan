package testdata

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func PresetUser() []*model.User {
	pwdFromWeb := utils.CalcSHA256("123456")

	return []*model.User{
		{
			UserName:  "admin",
			Nickname:  "admin",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			IsAdmin:   true,
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user",
			Nickname:  "user",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			Enable2FA: false,
			TotpKey:   "",
		},
		{
			UserName:  "user_with_totp",
			Nickname:  "user_with_totp",
			Password:  password.GeneratePwdHash(pwdFromWeb),
			Enable2FA: true,
			TotpKey:   "NVQXE2LP", // base32 of 'mario'
		},
	}
}
