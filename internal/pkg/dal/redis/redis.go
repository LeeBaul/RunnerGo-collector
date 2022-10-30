package redis

import (
	"fmt"
	"kp-collector/internal/pkg"
	"kp-collector/internal/pkg/dal/kao"
	"kp-collector/internal/pkg/log"
	"strconv"
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

func UpdatePartitionStatus(key string, partition int32) (err error) {
	field := fmt.Sprintf("%d", partition)
	err = RDB.HDel(key, field).Err()
	return
}

func InsertTestData(machineMap map[string]map[string]bool, sceneTestResultDataMsg *kao.SceneTestResultDataMsg) (err error) {
	data := sceneTestResultDataMsg.ToJson()
	key := fmt.Sprintf("%d:%s:reportData", sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId)
	if sceneTestResultDataMsg.End {
		reportId, err := strconv.ParseInt(sceneTestResultDataMsg.ReportId, 10, 64)
		if err != nil {
			log.Logger.Error("报告Id转数字失败：  ", err)
		}
		pkg.SendStopStressReport(machineMap, reportId)
	}

	err = RDB.LPush(key, data).Err()
	if err != nil {
		return
	}
	return
}

func Insert(rdb *redis.Client, a string) (err error) {
	err = rdb.LPush("report1", a).Err()
	if err != nil {
		return
	}
	return
}

type A struct {
	B int `json:"a"`
}
