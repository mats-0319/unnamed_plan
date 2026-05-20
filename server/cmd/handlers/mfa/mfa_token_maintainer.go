package mfa

import (
	"sync/atomic"
	"time"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/token"
)

type MFATokenMaintainer struct {
	Data map[string]*MFATokenItem // user name - mfa token item
}

type MFATokenItem struct {
	ExpireTime int64
	TryTimes   int
}

var mtm = &MFATokenMaintainer{Data: make(map[string]*MFATokenItem)}
var mtmLock atomic.Bool // 阈值全面清理的锁，保证同时只开一个goroutine执行清理

func newMFAToken(userName string, expireTime int64) {
	mtm.Data[userName] = &MFATokenItem{ExpireTime: expireTime}
}

func maintainMFAToken(tokenStr string) (t *token.Token, e *utils.Error) {
	go clearExpiredMFAToken()

	t, e = token.DeserializeToken(tokenStr, token.TokenType_MFAToken)
	if e != nil {
		return
	}

	v, ok := mtm.Data[t.UserName] // 能执行到这，t必然是不空的
	if !ok || v == nil || v.TryTimes >= 5 {
		delete(mtm.Data, t.UserName)
		return
	}

	mtm.Data[t.UserName].TryTimes++

	if v.ExpireTime < time.Now().UnixMilli() {
		delete(mtm.Data, t.UserName) // 访问到一个过期的token，删除它
		return
	}

	delete(mtm.Data, t.UserName) // 验证通过，token作废

	return
}

// clearExpiredMFAToken 积压token达到阈值时，全面检查并清理过期token
func clearExpiredMFAToken() {
	if len(mtm.Data) < 1000 {
		return
	}

	if mtmLock.CompareAndSwap(false, true) {
		defer mtmLock.Store(false)

		// 检查：因为只删除过期的token，所以即使在遍历map的过程中将新加入map的key也遍历出来了，也不会出错
		// （go不保证在这种情况下，新加入的key是否会被遍历）
		for k, v := range mtm.Data {
			if v != nil && v.ExpireTime < time.Now().UnixMilli() {
				delete(mtm.Data, k)
			}
		}
	}
}
