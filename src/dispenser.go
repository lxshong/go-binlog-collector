package src

import "errors"

// 分发器
type dispenser interface {
	Send(event *Event) error
}

// 分发器
func NewDispenser(instance string, config *InstanceConfig) (dispenser, error) {
	switch config.ToType {
	case REDIS:
		return newRedisDispenser(instance, config.Redis)
	}
	return nil, errors.New("to tyoe unsupport")
}
