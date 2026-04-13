package password

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"golang.org/x/crypto/argon2"
)

// pwd hash structure: `argon2id.v=19,m=65536,t=3,p=1.[saltHex].[keyHex]`
// more: doc/design.md 密码存储hash

type AlgorithmParams struct {
	CalcTimes uint32 // 迭代次数
	Memory    uint32 // 使用内存
	Threads   uint8  // 使用线程数
	KeyLength uint32
	Salt      []byte
}

func defaultAlgorithmParams() *AlgorithmParams {
	return &AlgorithmParams{
		CalcTimes: 3,
		Memory:    64 * 1024, // 64 MB
		Threads:   1,
		KeyLength: 32,
		Salt:      utils.GenerateRandomBytes[[]byte](32),
	}
}

func GeneratePassword(pwdSHA256 string) (pwdArgon2 string) {
	params := defaultAlgorithmParams()

	key := argon2.IDKey([]byte(pwdSHA256), params.Salt, params.CalcTimes, params.Memory, params.Threads, params.KeyLength)

	pwdArgon2 = fmt.Sprintf("argon2id.v=%d,m=%d,t=%d,c=%d.%s.%s", argon2.Version,
		params.Memory, params.CalcTimes, params.Threads, hex.EncodeToString(params.Salt), hex.EncodeToString(key))

	return
}

// VerifyPassword decode 'key' from 'pwdArgon2', calc 'new key' with 'pwdSHA256', compare two keys
func VerifyPassword(pwdSHA256 string, pwdArgon2 string) *utils.Error {
	params, oldKey, e := decodeHash(pwdArgon2)
	if e != nil {
		return e
	}

	newKey := argon2.IDKey([]byte(pwdSHA256), params.Salt, params.CalcTimes, params.Memory, params.Threads, params.KeyLength)

	if subtle.ConstantTimeCompare(oldKey, newKey) != 1 { // 使用恒定时间比较防止时序攻击
		e := utils.ErrWrongPassword().WithParam("old key", oldKey).WithParam("new key", newKey)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func decodeHash(pwdHash string) (params *AlgorithmParams, oldKey []byte, e *utils.Error) {
	pwdSplit := strings.Split(pwdHash, ".")
	if len(pwdSplit) != 4 || pwdSplit[0] != "argon2id" {
		e = utils.ErrInvalidPassword().WithParam("encoded pwd", pwdHash)
		mlog.Error(e.String())
		return
	}

	var version int
	params = &AlgorithmParams{}
	_, err := fmt.Sscanf(pwdSplit[1], "v=%d,m=%d,t=%d,c=%d", &version, &params.Memory, &params.CalcTimes, &params.Threads)
	if err != nil || version != argon2.Version {
		e = utils.ErrInvalidPassword().WithCause(err).WithParam("version", version).WithParam("params", params)
		mlog.Error(e.String())
		return
	}

	params.Salt, err = hex.DecodeString(pwdSplit[2])
	if err != nil {
		e = utils.ErrInvalidPwdSalt().WithCause(err).WithParam("salt", pwdSplit[2])
		mlog.Error(e.String())
		return
	}

	key, err := hex.DecodeString(pwdSplit[3])
	if err != nil {
		e = utils.ErrInvalidPwdKey().WithCause(err).WithParam("key", pwdSplit[3])
		mlog.Error(e.String())
		return
	}
	params.KeyLength = uint32(len(key))

	return
}
