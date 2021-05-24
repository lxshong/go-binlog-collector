package src

import (
	"errors"
)

type collector interface {
	Run(call func(event *Event) error) error
}

// 创建收集器
func NewCollector(instance string, config *InstanceConfig) (collector, error) {
	switch config.FromType {
	case MYSQL:
		return newMysqlCollector(instance, config.Mysql)
	}
	return nil, errors.New("instance type unsupport")
}
