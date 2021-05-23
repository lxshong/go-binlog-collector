package src

import "errors"

// 分发器
type dispenser interface {
	Send(event *Event) error
}

// 分发器
func NewDispenser(instance string, toType string, config interface{}) (dispenser, error) {
	switch toType {
	case REDIS:
		return newRedisDispenser(instance, config.(*RedisConfig))
	}
	return nil, errors.New("to tyoe unsupport")
}
