package utils

import (
	"testing"
)

func TestPassword(t *testing.T) {
	pwdFromWeb := CalcSHA256("123456")
	t.Log("> Password From Web: ", pwdFromWeb)

	pwdDB := GeneratePwdHash(pwdFromWeb)
	t.Log("> Password To DB: ", pwdDB)

	err := VerifyPassword(pwdFromWeb, pwdDB)
	if err != nil {
		t.Fatal(err.String())
	}
	t.Log("> Verified.")
}
