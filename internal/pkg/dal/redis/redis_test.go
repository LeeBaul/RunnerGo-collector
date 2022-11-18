package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestInsert(t *testing.T) {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "",
			Password: "",
			DB:       0,
		})

	var a = &A{}
	for i := 0; i < 10; i++ {
		a.B = i
		s, _ := json.Marshal(a)
		err := Insert(rdb, string(s))
		if err != nil {
			fmt.Println(err)
		}
	}

	val := rdb.LRange("report1", 0, -1).Val()
	for i := len(val) - 1; i >= 0; i-- {
		fmt.Println("result:         ", val[i])
	}

}
