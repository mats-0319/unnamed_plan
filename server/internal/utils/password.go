package utils

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"golang.org/x/crypto/argon2"
)

// hash str structure: `argon2id.v=19.m=65536,t=3,p=1.[saltHex].[keyHex]`

type PasswordManager struct {
	CalcTimes uint32 // 迭代次数
	Memory    uint32 // 使用内存
	Threads   uint8  // 使用线程数
	KeyLength uint32
}

func defaultPasswordManager() *PasswordManager {
	return &PasswordManager{
		CalcTimes: 3,
		Memory:    64 * 1024, // 64 MB
		Threads:   1,
		KeyLength: 32,
	}
}

func GeneratePwdHash(password string, params ...*PasswordManager) string {
	var pm = defaultPasswordManager()
	if len(params) > 1 {
		pm = params[0]
	}

	salt := GenerateRandomBytes[[]byte](32)
	saltHex := hex.EncodeToString(salt)

	key := argon2.IDKey([]byte(password), salt, pm.CalcTimes, pm.Memory, pm.Threads, pm.KeyLength)
	keyHex := hex.EncodeToString(key)

	return fmt.Sprintf("argon2id.v=19.m=%d,t=%d,p=%d.%s.%s",
		pm.Memory, pm.CalcTimes, pm.Threads, saltHex, keyHex)
}

func VerifyPassword(password, encodedHash string) *Error {
	params, salt, key, err := decodeHash(encodedHash)
	if err != nil {
		return err
	}

	newKey := argon2.IDKey([]byte(password), salt, params.CalcTimes, params.Memory, params.Threads, params.KeyLength)

	// 使用恒定时间比较防止时序攻击
	if subtle.ConstantTimeCompare(key, newKey) != 1 {
		e := NewError(ET_ParamsError, ED_WrongPwd).WithParam("old key", key).WithParam("new key", newKey)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func decodeHash(encodedHash string) (*PasswordManager, []byte, []byte, *Error) {
	pwdSplit := strings.Split(encodedHash, ".")
	if len(pwdSplit) != 5 || pwdSplit[0] != "argon2id" {
		e := NewError(ET_ParamsError, ED_PwdStructure).WithParam("encoded pwd", encodedHash)
		mlog.Log(e.String())
		return nil, nil, nil, e
	}

	var version int
	_, err := fmt.Sscanf(pwdSplit[1], "v=%d", &version)
	if err != nil || version != argon2.Version {
		e := NewError(ET_ParamsError, ED_PwdVersion).WithCause(err).WithParam("version", version)
		mlog.Log(e.String())
		return nil, nil, nil, e
	}

	params := &PasswordManager{}
	_, err = fmt.Sscanf(pwdSplit[2], "m=%d,t=%d,p=%d", &params.Memory, &params.CalcTimes, &params.Threads)
	if err != nil {
		e := NewError(ET_ParamsError, ED_PwdParams).WithCause(err)
		mlog.Log(e.String())
		return nil, nil, nil, e
	}

	salt, err := hex.DecodeString(pwdSplit[3])
	if err != nil {
		e := NewError(ET_ParamsError, ED_HexDecode).WithCause(err).WithParam("salt", pwdSplit[3])
		mlog.Log(e.String())
		return nil, nil, nil, e
	}

	key, err := hex.DecodeString(pwdSplit[4])
	if err != nil {
		e := NewError(ET_ParamsError, ED_HexDecode).WithCause(err).WithParam("key", pwdSplit[4])
		mlog.Log(e.String())
		return nil, nil, nil, e
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
