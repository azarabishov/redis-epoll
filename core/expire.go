package core

import (
	"time"
)

func hasExpired(obj *Obj) bool {
	exp, ok := expires[obj]
	if !ok {
		return false
	}
	return exp <= uint64(time.Now().UnixMilli())
}

func getExpiry(obj *Obj) (uint64, bool) {
	exp, ok := expires[obj]
	return exp, ok
}

func expireSample() float32 {
	var limit int = 20
	var expiredCount int = 0

	for keyPtr, obj := range store {
		limit--
		if hasExpired(obj) {
			DelByPtr(keyPtr)
			expiredCount++
		}

		if limit == 0 {
			break
		}
	}

	return float32(expiredCount) / float32(20.0)
}

func DeleteExpiredKeys() {
	for {
		frac := expireSample()
		if frac < 0.25 {
			break
		}
	}
}
