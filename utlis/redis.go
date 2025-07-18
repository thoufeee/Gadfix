package utlis

import (
	"fmt"
	"gadfix/db"

	"time"
)

// storing refresh token
func StoreRefresh(tokenid string, userid uint, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", tokenid)
	return db.Redis.Set(db.Ctx, key, userid, expiry).Err()
}

// delete refresh token
func DeleteRefresh(tokenid string) error {
	key := fmt.Sprintf("refresh_token:%s", tokenid)
	return db.Redis.Del(db.Ctx, key).Err()
}

// valid check refresh token
func ValidRefresh(tokenid string) bool {
	key := fmt.Sprintf("refresh_token:%s", tokenid)
	_, err := db.Redis.Get(db.Ctx, key).Result()
	return err == nil
}
