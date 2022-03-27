package convert

import (
	"context"
	"errors"
	"go-binlog-collector/src/utils"

	"github.com/go-mysql-org/go-mysql/replication"
)

type Converter interface {
	Do(ctx context.Context, binlogEvent *replication.BinlogEvent) (*utils.Event, error)
}

// 创建收集器
func NewConverter(config *utils.InstanceConfig) (Converter, error) {
	switch config.FromType {
	case utils.MYSQL:
		return newMysqlConverter(config.Instance, config.Mysql, config.Rules), nil
	}
	return nil, errors.New("instance type unsupport")
}
