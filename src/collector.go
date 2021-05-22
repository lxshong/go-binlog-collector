package src

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type collector interface {
	Run() error
}

type event struct {
	EventType string
	Schema    string
	Table     string
}

func NewCollector(instance string) (collector, error) {
	config, err := getInstanceConfig(instance)
	if err != nil {
		return nil, err
	}
	switch config.InstanceType {
	case "mysql":
		return &mysqlCollector{
			flavor:   config.InstanceType,
			instance: instance,
			config:   config.Config.(MysqlConfig),
			Pos:      NewPosition(instance),
		}, nil
	}
	return nil, errors.New("instance type unsupport")
}

// 收集器
type mysqlCollector struct {
	flavor   string
	instance string
	config   MysqlConfig
	Pos      *position
}

func (collector *mysqlCollector) Run() error {

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
		} else if event != nil{
			fmt.Println(event)
		}
		pos = syncer.GetNextPosition()
		collector.Pos.Name = pos.Name
		collector.Pos.Pos = int(pos.Pos)
		if err := collector.Pos.Save(); err != nil {
			return err
		}
	}
}
