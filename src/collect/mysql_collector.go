package collect

import (
	"context"
	"go-binlog-collector/src/convert"
	"go-binlog-collector/src/utils"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

func newMysqlCollector(config *utils.InstanceConfig) (Collector, error) {
	converter, err := convert.NewConverter(config)
	if err != nil {
		return nil, err
	}
	return &mysqlCollector{
		flavor:    utils.MYSQL,
		instance:  config.Instance,
		config:    config.Mysql,
		converter: converter,
		Pos:       utils.NewPosition(config.Instance),
	}, nil
}

// 收集器
type mysqlCollector struct {
	flavor    string
	instance  string
	config    *utils.MysqlConfig
	converter convert.Converter
	Pos       *utils.Position
}

func (collector *mysqlCollector) Do(ctx context.Context, call func(event *utils.Event) error) error {

	cfg := replication.BinlogSyncerConfig{
		ServerID: 101,
		Flavor:   collector.flavor,
		Host:     collector.config.Host,
		Port:     uint16(collector.config.Port),
		User:     collector.config.User,
		Password: collector.config.Passwd,
	}
	if err := collector.Pos.Init(); err != nil {
		return err
	}
	syncer := replication.NewBinlogSyncer(cfg)
	// Start sync with specified binlog file and position
	pos := mysql.Position{collector.Pos.Name, uint32(collector.Pos.Pos)}
	streamer, _ := syncer.StartSync(pos)
	for {
		ev, _ := streamer.GetEvent(context.Background())
		if ev == nil {
			continue
		}
		// 事件转换
		if event, err := collector.converter.Do(ctx, ev); err != nil {
			return err
		} else {
			if event != nil {
				// 分发
				if err := call(event); err != nil {
					return err
				}
			}
			// 记录位置
			pos = syncer.GetNextPosition()
			collector.Pos.Name = pos.Name
			collector.Pos.Pos = int(pos.Pos)
			if err := collector.Pos.Save(); err != nil {
				return err
			}
		}
	}
}
