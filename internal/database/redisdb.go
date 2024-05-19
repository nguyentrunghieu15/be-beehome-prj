package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	*redis.Client
}

func (*RedisDb) Init() interface{} {
	addr := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	for {
		result := client.Ping(context.Background())
		if result.Err() != nil {
			log.Println(result.Err())
			time.Sleep(300 * time.Millisecond)
			continue
		}
		break
	}
	return &RedisDb{
		Client: client,
	}
}
