package redis

import (
	"fmt"
	"kp-collector/internal/pkg/dal/kao"
	"time"
)
import "github.com/go-redis/redis"

var (
	RDB          *redis.Client
	timeDuration = 3 * time.Second
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(addr, password string, db int64) (err error) {
	RDB = redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		})
	_, err = RDB.Ping().Result()
	return err
}

func DeletePartition(key string) (err error) {
	_, err = RDB.Get(key).Result()
	if err != nil {
		return
	}
	_, err = RDB.Del(key).Result()
	if err != nil {
		return
	}
	return
}

func InsertTestData(sceneTestResultDataMsg kao.SceneTestResultDataMsg) (err error) {
	data := sceneTestResultDataMsg.ToJson()
	key := fmt.Sprintf("%d:%s:reportData", sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId)
	err = RDB.LPush(key, data).Err()
	if err != nil {
		return
	}
	return
}

func Insert(rdb *redis.Client, a *A) (err error) {
	err = rdb.LPush("report1", a).Err()
	if err != nil {
		return
	}
	return
}

type A struct {
	B int `json:"a"`
}
