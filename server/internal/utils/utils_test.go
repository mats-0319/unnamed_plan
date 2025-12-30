package utils

import "testing"

func TestPasswordWithSalt(t *testing.T) {
	salt := string(GenerateRandomBytes(10))
	password := "123456"
	password = CalcSHA256(password)
	t.Log("pwd 1:", password)

	password = CalcSHA256(password + salt)
	t.Log("pwd 2:", password)
	t.Log("salt :", salt)

	//    utils_test.go:9: pwd 1: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
	//    utils_test.go:12: pwd 2: 5a6e00fc37b82fabdd017b978928f4e79a7d8decfcbc67a46ea8a9f60fdeb85a
	//    utils_test.go:13: salt : hCvQqZwe1L
}
