package src

import "errors"

// 分发器
type dispenser interface {
	Send(event *Event) error
}

func NewDispenser(instance string) (dispenser, error) {
	config, err := getInstanceConfig(instance)
	if err != nil {
		return nil, err
	}
	switch config.ToType {
	case "redis":
		return newRedisDispenser(instance)
	}
	return nil, errors.New("to tyoe unsupport")
}
