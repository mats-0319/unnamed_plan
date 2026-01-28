package utils

import "testing"

func TestPasswordWithSalt(t *testing.T) {
	salt := GenerateRandomBytes[string](10)
	t.Log("salt :", salt)

	password := "123456"
	password = HmacSHA256[string](password)
	t.Log("pwd 1:", password)

	password = HmacSHA256(password, salt)
	t.Log("pwd 2:", password)

	//    utils_test.go:7:  salt : 3RHPTqCpQ7
	//    utils_test.go:11: pwd 1: dbd978ccdbbe8b6de77f6b37b5df9b5b62a7e892a501c3b53eaa16b0838bd5ed
	//    utils_test.go:14: pwd 2: d0cdec83c70294c327df7808118fdf80697438664361b69095496075d1d616e1
}
