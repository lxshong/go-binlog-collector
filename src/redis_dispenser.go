package src

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

// 创建RedisDispenser
func newRedisDispenser(instance string) (dispenser, error) {
	// 读取配置
	config, err := getInstanceConfig(instance)
	if err != nil {
		return nil, err
	}
	redisConfig, ok := config.ToConfig.(*RedisConfig)
	if !ok {
		return nil, errors.New("redis config convert failed")
	}
	// redis 是否可用
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.Passwd,
		DB:       redisConfig.DB,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	// 创建dispenser
	if config.ToType == "redis" {
		return &redisDispenser{
			prekey: instance,
			config: redisConfig,
			client: client,
		}, nil
	}
	return nil, errors.New("redis config do not match")
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
