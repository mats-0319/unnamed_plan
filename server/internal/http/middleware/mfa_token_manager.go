package middleware

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync/atomic"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type MfaTokenManager struct {
	Data map[string]*MfaTokenItem // user name - mfa token item
}

type MfaTokenItem struct {
	Token      string // complete token
	ExpireTime int64
	TryTimes   int
}

var mtm = &MfaTokenManager{Data: make(map[string]*MfaTokenItem)}
var clearLock atomic.Bool // 阈值全面清理的锁，保证同时只开一个goroutine执行清理

func NewMfaToken(userName string, token string, expireTime int64) {
	mtm.Data[userName] = &MfaTokenItem{Token: token, ExpireTime: expireTime}
}

func VerifyMfaToken(userName string, token string) (e *utils.Error) {
	clearExpiredMfaToken()

	v, ok := mtm.Data[userName]
	if !ok || v == nil || v.TryTimes >= 5 {
		if v == nil || v.TryTimes >= 5 {
			delete(mtm.Data, userName)
		}

		e = utils.ErrNoMfaToken().WithParam("user name", userName)
		return
	}
	mtm.Data[userName].TryTimes++

	if v.Token != token { // 只验证token一致，不展开比较具体属性
		e = utils.ErrWrongMfaToken().WithParam("token in memory", v).WithParam("token", token)
		return
	}

	if v.ExpireTime < time.Now().UnixMilli() {
		delete(mtm.Data, userName) // 访问时删除，参考redis的过期策略

		e = utils.ErrMfaTokenExpired().WithParam("expire time", v.ExpireTime)
		return
	}

	delete(mtm.Data, userName) // 验证通过，token作废

	return
}

func DecodeUserNameFromMfaToken(token string) (userName string, e *utils.Error) {
	tokenSplit := strings.Split(token, ".")
	if len(tokenSplit) != 2 {
		e = utils.ErrInvalidMfaToken().WithParam("token", token)
		mlog.Error(e.String())
		return
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = utils.ErrDecodeMfaToken().WithCause(err)
		mlog.Error(e.String())
		return
	}

	tokenIns := &Token{}
	if err := json.Unmarshal(tokenBytes, tokenIns); err != nil {
		e = utils.ErrDeserializeMfaToken().WithCause(err)
		mlog.Error(e.String())
		return
	}

	if tokenIns.Type != TokenType_MfaToken { // 其实这里不用验证，后面会验证传入token和内存中的值是否一致
		e = utils.ErrInvalidMfaToken().WithParam("token type", tokenIns.Type)
		mlog.Error(e.String())
		return
	}

	userName = tokenIns.UserName

	return
}

// clearExpiredMfaToken 积压token达到阈值时，全面检查并清理过期token
func clearExpiredMfaToken() {
	if len(mtm.Data) > 1000 && clearLock.CompareAndSwap(false, true) {
		go func() {
			defer clearLock.Store(false)

			for k, v := range mtm.Data {
				if v != nil && v.ExpireTime < time.Now().UnixMilli() {
					delete(mtm.Data, k)
				}
			}
		}()
	}
}
