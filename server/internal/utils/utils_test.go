package utils

import "testing"

func TestPasswordWithSalt(t *testing.T) {
	salt := GenerateRandomStr(10)
	t.Log("salt :", salt)

	password := "123456"
	password = CalcSHA256(password)
	t.Log("pwd 1:", password)

	password = CalcSHA256(password + salt)
	t.Log("pwd 2:", password)

	//    utils_test.go:7: salt : vkAvvZNxfO
	//    utils_test.go:11: pwd 1: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
	//    utils_test.go:14: pwd 2: 0fae47e6f0d24e8366f55484e5af228f5eefeed3ca7e5cb2f1f9d31b3579bae5
}
