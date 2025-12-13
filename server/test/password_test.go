package test

import (
	"testing"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func TestPasswordWithSalt(t *testing.T) {
	salt := string(utils.GenerateRandomBytes(10))
	password := "123456"
	password = utils.CalcSHA256(password)
	t.Log("pwd 1:", password)

	password = utils.CalcSHA256(password + salt)
	t.Log("pwd 2:", password)
	t.Log("salt :", salt)
}
