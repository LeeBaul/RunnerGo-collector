package redis

import (
	"RunnerGo-collector/internal/pkg"
	"RunnerGo-collector/internal/pkg/dal/kao"
	"RunnerGo-collector/internal/pkg/log"
	"fmt"
	"strconv"
	"time"
)
import "github.com/go-redis/redis"

var (
	RDB1         *redis.Client
	RDB2         *redis.Client
	timeDuration = 3 * time.Second
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(addr1, password1 string, db1 int64, addr2, password2 string, db2 int64) (err error) {
	RDB1 = redis.NewClient(
		&redis.Options{
			Addr:     addr1,
			Password: password1,
			DB:       int(db1),
		})
	_, err = RDB1.Ping().Result()
	if err != nil {
		return
	}
	RDB2 = redis.NewClient(
		&redis.Options{
			Addr:     addr2,
			Password: password2,
			DB:       int(db2),
		})
	_, err = RDB2.Ping().Result()
	return err
}

func UpdatePartitionStatus(key string, partition int32) (err error) {
	field := fmt.Sprintf("%d", partition)
	err = RDB2.HDel(key, field).Err()
	return
}

func InsertTestData(machineMap map[string]map[string]int64, sceneTestResultDataMsg *kao.SceneTestResultDataMsg) (err error) {
	data := sceneTestResultDataMsg.ToJson()
	key := fmt.Sprintf("%d:%s:reportData", sceneTestResultDataMsg.PlanId, sceneTestResultDataMsg.ReportId)
	if sceneTestResultDataMsg.End {
		reportId, err := strconv.ParseInt(sceneTestResultDataMsg.ReportId, 10, 64)
		if err != nil {
			log.Logger.Error("报告Id转数字失败：  ", err)
		}
		pkg.SendStopStressReport(machineMap, reportId)
	}

	err = RDB1.LPush(key, data).Err()
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
