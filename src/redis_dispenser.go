package src

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

// 创建RedisDispenser
func newRedisDispenser(instance string, config *RedisConfig) (dispenser, error) {
	// redis 是否可用
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Passwd,
		DB:       config.DB,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	// 创建dispenser
	return &redisDispenser{
		prekey: instance,
		config: config,
		client: client,
	}, nil
}

type redisDispenser struct {
	prekey string
	config *RedisConfig
	client *redis.Client
}

// 发送
func (d *redisDispenser) Send(event *Event) error {
	if d.config == nil {
		return errors.New("redis config unexists")
	}

	//连接服务器
	if d.client == nil {

	}
	key := fmt.Sprintf("%s:%s:%s", d.prekey, event.DataBase, event.Table)
	eventStr := event.String()
	if err := d.client.RPush(key, eventStr).Err(); err != nil {
		return err
	}
	fmt.Println("redis send success:", eventStr)
	return nil
}
