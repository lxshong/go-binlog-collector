package src

import (
	"errors"
)

type collector interface {
	Run(call func(event *Event) error) error
}

// 创建收集器
func NewCollector(instance string, fromType string, config interface{}) (collector, error) {
	switch fromType {
	case MYSQL:
		return newMysqlCollector(instance, config.(*MysqlConfig))
	}
	return nil, errors.New("instance type unsupport")
}
