package redis

import "time"
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
