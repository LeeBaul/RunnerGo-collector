package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestInsert(t *testing.T) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "172.17.101.191:6398",
			Password: "apipost",
			DB:       0,
		})

	var a = &A{}
	for i := 0; i < 10; i++ {
		a.B = i
		err := Insert(rdb, a)
		if err != nil {
			fmt.Println(err)
		}
	}

	val := rdb.LRange("report1", 0, -1).Val()
	fmt.Println("result:         ", val)
}
