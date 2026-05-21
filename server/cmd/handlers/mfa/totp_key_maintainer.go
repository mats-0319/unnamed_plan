package mfa

import (
	"sync/atomic"
	"time"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type TOTPKeyMaintainer struct {
	Data map[string]*TOTPKeyItem // username - totp key item
	lock atomic.Bool
}

type TOTPKeyItem struct {
	Key        string // origin length: 10, base32 encoded (length: 16)
	ExpireTime int64
}

var tkm = &TOTPKeyMaintainer{Data: make(map[string]*TOTPKeyItem)}

func newTOTPKey(userName string, key string, expireTime int64) {
	tkm.Data[userName] = &TOTPKeyItem{Key: key, ExpireTime: expireTime}
}

func maintainTOTPKey(userName string, code string) (key string, e *utils.Error) {
	go clearExpiredTOTPKey()

	v, ok := tkm.Data[userName]
	if !ok {
		return
	}

	if v.ExpireTime < time.Now().UnixMilli() {
		e = utils.ErrTokenExpired()
		delete(tkm.Data, userName) // 访问的缓存已过期，删除
		return
	}

	if e = VerifyTOTPCode(code, v.Key); e != nil {
		return
	}

	delete(tkm.Data, userName) // 缓存使用成功

	key = v.Key

	return
}

func clearExpiredTOTPKey() {
	if len(tkm.Data) < 1000 {
		return
	}

	if tkm.lock.CompareAndSwap(false, true) {
		defer tkm.lock.Store(false)

		// 检查：因为只删除过期的token，所以即使在遍历map的过程中将新加入map的key也遍历出来了也不会出错
		// （go不保证在这种情况下，新加入的key是否会被遍历）
		for k, v := range tkm.Data {
			if v != nil && v.ExpireTime < time.Now().UnixMilli() {
				delete(tkm.Data, k)
			}
		}
	}
}
