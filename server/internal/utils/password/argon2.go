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
}

func defaultAlgorithmParams() *AlgorithmParams {
	return &AlgorithmParams{
		CalcTimes: 3,
		Memory:    64 * 1024, // 64 MB
		Threads:   1,
		KeyLength: 32,
	}
}

func GeneratePwdHash(password string, params ...*AlgorithmParams) string {
	var pm = defaultAlgorithmParams()
	if len(params) > 1 {
		pm = params[0]
	}

	salt := utils.GenerateRandomBytes[[]byte](32)
	saltHex := hex.EncodeToString(salt)

	key := argon2.IDKey([]byte(password), salt, pm.CalcTimes, pm.Memory, pm.Threads, pm.KeyLength)
	keyHex := hex.EncodeToString(key)

	return fmt.Sprintf("argon2id.v=%d,m=%d,t=%d,c=%d.%s.%s",
		argon2.Version, pm.Memory, pm.CalcTimes, pm.Threads, saltHex, keyHex)
}

// VerifyPassword decode 'key' from 'pwd hash', calc 'new key' with 'password', compare two keys
func VerifyPassword(password, pwdHash string) *utils.Error {
	params, salt, key, err := decodeHash(pwdHash)
	if err != nil {
		return err
	}

	newKey := argon2.IDKey([]byte(password), salt, params.CalcTimes, params.Memory, params.Threads, params.KeyLength)

	// 使用恒定时间比较防止时序攻击
	if subtle.ConstantTimeCompare(key, newKey) != 1 {
		e := utils.ErrWrongPwd().WithParam("old key", key).WithParam("new key", newKey)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func decodeHash(pwdHash string) (*AlgorithmParams, []byte, []byte, *utils.Error) {
	pwdSplit := strings.Split(pwdHash, ".")
	if len(pwdSplit) != 4 || pwdSplit[0] != "argon2id" {
		e := utils.ErrInvalidPwd().WithParam("encoded pwd", pwdHash)
		mlog.Error(e.String())
		return nil, nil, nil, e
	}

	var version int
	params := &AlgorithmParams{}
	_, err := fmt.Sscanf(pwdSplit[1], "v=%d,m=%d,t=%d,c=%d", &version, &params.Memory, &params.CalcTimes, &params.Threads)
	if err != nil || version != argon2.Version {
		e := utils.ErrPwdParams().WithCause(err).WithParam("version", version).WithParam("params", params)
		mlog.Error(e.String())
		return nil, nil, nil, e
	}

	salt, err := hex.DecodeString(pwdSplit[2])
	if err != nil {
		e := utils.ErrHexDecode().WithCause(err).WithParam("salt", pwdSplit[2])
		mlog.Error(e.String())
		return nil, nil, nil, e
	}

	key, err := hex.DecodeString(pwdSplit[3])
	if err != nil {
		e := utils.ErrHexDecode().WithCause(err).WithParam("key", pwdSplit[3])
		mlog.Error(e.String())
		return nil, nil, nil, e
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
