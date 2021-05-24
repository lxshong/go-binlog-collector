package src

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

func newMysqlCollector(instance string, config *MysqlConfig, rules map[string]map[string]int) (collector, error) {
	converter, err := NewConverter(instance, MYSQL, config, rules)
	if err != nil {
		return nil, err
	}
	return &mysqlCollector{
		flavor:    MYSQL,
		instance:  instance,
		config:    config,
		converter: converter,
		Pos:       NewPosition(instance),
	}, nil
}

// 收集器
type mysqlCollector struct {
	flavor    string
	instance  string
	config    *MysqlConfig
	converter converter
	Pos       *position
}

func (collector *mysqlCollector) Run(call func(event *Event) error) error {

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
		if event, err := collector.converter.Convert(ev); err != nil {
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
