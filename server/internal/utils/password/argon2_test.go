package password

import (
	"testing"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func TestPassword(t *testing.T) {
	pwdFromWeb := utils.CalcSHA256("123456") // hex(hash('password'))
	t.Log("> Password From Web: ", pwdFromWeb)

	pwdDB := GeneratePwdHash(pwdFromWeb)
	t.Log("> Password To DB: ", pwdDB)

	err := VerifyPassword(pwdFromWeb, pwdDB)
	if err != nil {
		t.Fatal(err.String())
	}
	t.Log("> Verified.")
}
