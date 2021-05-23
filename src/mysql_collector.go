package src

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

func newMysqlCollector(instance string, config *MysqlConfig) (collector, error) {
	return &mysqlCollector{
		flavor:   MYSQL,
		instance: instance,
		config:   config,
		Pos:      NewPosition(instance),
	}, nil
}

// 收集器
type mysqlCollector struct {
	flavor   string
	instance string
	config   *MysqlConfig
	Pos      *position
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
	converter, err := NewConverter(collector.instance)
	if err != nil {
		return err
	}
	syncer := replication.NewBinlogSyncer(cfg)
	// Start sync with specified binlog file and position
	pos := mysql.Position{collector.Pos.Name, uint32(collector.Pos.Pos)}
	streamer, _ := syncer.StartSync(pos)
	for {
		ev, _ := streamer.GetEvent(context.Background())
		if event, err := converter.Convert(ev); err != nil {
			return err
		} else {
			if event != nil {
				if err := call(event); err != nil {
					return err
				}
			}
			pos = syncer.GetNextPosition()
			collector.Pos.Name = pos.Name
			collector.Pos.Pos = int(pos.Pos)
			if err := collector.Pos.Save(); err != nil {
				return err
			}
		}
	}
}
