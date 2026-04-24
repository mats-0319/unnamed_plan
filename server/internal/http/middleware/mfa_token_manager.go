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

type MFATokenManager struct {
	Data map[string]*MFATokenItem // user name - mfa token item
}

type MFATokenItem struct {
	Token      string // complete token
	ExpireTime int64
	TryTimes   int
}

var mtm = &MFATokenManager{Data: make(map[string]*MFATokenItem)}
var clearLock atomic.Bool // 阈值全面清理的锁，保证同时只开一个goroutine执行清理

func NewMFAToken(userName string, token string, expireTime int64) {
	mtm.Data[userName] = &MFATokenItem{Token: token, ExpireTime: expireTime}
}

func VerifyMFAToken(userName string, token string) (e *utils.Error) {
	clearExpiredMFAToken()

	v, ok := mtm.Data[userName]
	if !ok || v == nil || v.TryTimes >= 5 {
		if v == nil || v.TryTimes >= 5 {
			delete(mtm.Data, userName)
		}

		e = utils.ErrNoMFAToken().WithParam("user name", userName)
		mlog.Error(e.String())
		return
	}
	mtm.Data[userName].TryTimes++

	if v.Token != token { // 只验证token是否一致，不展开比较具体属性
		e = utils.ErrWrongMFAToken().WithParam("token in memory", v).WithParam("token", token)
		mlog.Error(e.String())
		return
	}

	if v.ExpireTime < time.Now().UnixMilli() {
		delete(mtm.Data, userName) // 访问时删除

		e = utils.ErrMFATokenExpired().WithParam("expire time", v.ExpireTime)
		mlog.Error(e.String())
		return
	}

	delete(mtm.Data, userName) // 验证通过，token作废

	return
}

func DecodeUserNameFromMFAToken(token string) (userName string, e *utils.Error) {
	tokenSplit := strings.Split(token, ".")
	if len(tokenSplit) != 2 {
		e = utils.ErrInvalidMFAToken().WithParam("token", token)
		mlog.Error(e.String())
		return
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = utils.ErrDecodeMFAToken().WithCause(err)
		mlog.Error(e.String())
		return
	}

	tokenIns := &Token{}
	if err := json.Unmarshal(tokenBytes, tokenIns); err != nil {
		e = utils.ErrDeserializeMFAToken().WithCause(err)
		mlog.Error(e.String())
		return
	}

	if tokenIns.Type != TokenType_MFAToken { // 其实这里不用验证，后面会验证传入token和内存中的值是否一致
		e = utils.ErrInvalidMFAToken().WithParam("token type", tokenIns.Type)
		mlog.Error(e.String())
		return
	}

	userName = tokenIns.UserName

	return
}

// clearExpiredMFAToken 积压token达到阈值时，全面检查并清理过期token
func clearExpiredMFAToken() {
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
