package dispense

import (
	"errors"
	"go-binlog-collector/src/utils"
)

// 分发器
type dispenser interface {
	Do(event *utils.Event) error
}

// 分发器
func NewDispenser(config *utils.InstanceConfig) (dispenser, error) {
	switch config.ToType {
	case utils.STDIO:
		return newStdioDispenser(config.Instance, config.Redis)
	case utils.REDIS:
		return newRedisDispenser(config.Instance, config.Redis)
	}
	return nil, errors.New("to tyoe unsupport")
}
