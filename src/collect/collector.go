package collect

import (
	"context"
	"errors"
	"go-binlog-collector/src/utils"
)

type Collector interface {
	Do(ctx context.Context, call func(event *utils.Event) error) error
}

// 创建收集器
func NewCollector(config *utils.InstanceConfig) (Collector, error) {
	switch config.FromType {
	case utils.MYSQL:
		return newMysqlCollector(config)
	}
	return nil, errors.New("instance type unsupport")
}
