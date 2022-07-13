package common

import (
	"log"
	"testing"
	"time"
)

func TestRedisConnect(t *testing.T) {
	TestSetup(t)
	InitRedisClient()
	redisClient := RedisClient()
	err := redisClient.Set(ctx, "test", "123", time.Minute*10).Err()
	if err != nil {
		t.Error(err)
	}
	stringCmd := redisClient.Get(ctx, "test")
	if stringCmd.Err() != nil {
		t.Error(stringCmd.Err())
	}
	log.Println("get value is :" + stringCmd.Val())
	if stringCmd.Val() != "123" {
		t.Error("redis get value error")
	}

}
