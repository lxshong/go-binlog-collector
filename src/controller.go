package src

import (
	"context"
	"go-binlog-collector/src/collect"
	"go-binlog-collector/src/dispense"
	"go-binlog-collector/src/utils"
)

// 执行
func Do(ctx context.Context, confPath string) error {
	config, err := utils.ParseConfig(confPath)
	if err != nil {
		return err
	}
	collector, err := collect.NewCollector(config)
	if err != nil {
		return err
	}
	dispenser, err := dispense.NewDispenser(config)
	if err != nil {
		return err
	}
	return collector.Do(ctx, func(event *utils.Event) error {
		return dispenser.Do(event)
	})
}
