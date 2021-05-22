package src

import (
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
)

type Event struct {
	Table    string
	DataBase string
	Before   map[string]interface{}
	After    map[string]interface{}
}

// 创建转换器
func NewConverter(instance string) (converter, error) {
	config, err := getInstanceConfig(instance)
	if err != nil {
		return nil, err
	}
	switch config.InstanceType {
	case "mysql":
		cvt := new(mysqlConverter)
		cvt.instance = instance
		cvt.config = config.Config.(MysqlConfig)
		cvt.tables = make(map[string]schema.Table)
		return cvt, nil
	}
	return nil, errors.New("instance type unsupport")
}

type converter interface {
	Convert(binlogEvent *replication.BinlogEvent) (*Event, error)
}

type mysqlConverter struct {
	instance string
	tables   map[string]schema.Table
	config   MysqlConfig
}

// 转换
func (c *mysqlConverter) Convert(ev *replication.BinlogEvent) (*Event, error) {
	switch ev.Event.(type) {
	case *replication.RowsEvent:
		rowsEvent := ev.Event.(*replication.RowsEvent)
		if event, err := c.rowsEventConvert(rowsEvent); err != nil {
			return nil, err
		} else {
			return event, nil
		}
	}
	return nil, nil
}

func (c *mysqlConverter) GetTable(database string, table string) error {
	key := c.key(database, table)
	config :=  c.config
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, _ := client.Connect(addr, config.User, config.Passwd, config.DataBase)
	tableInfo, err := schema.NewTable(conn, config.DataBase, config.Table)
	if err != nil {
		return err
	}
	c.tables[key] = *tableInfo
	return nil
}

// RowsEvent 事件处理
func (c *mysqlConverter) rowsEventConvert(rowsEvent *replication.RowsEvent) (*Event, error) {
	// 获取表结构
	table := string(rowsEvent.Table.Table)
	database := string(rowsEvent.Table.Schema)
	if err := c.GetTable(database, table); err != nil {
		return nil, err
	}
	// 初始化event对象
	tableKey := c.key(database, table)
	tableInfo := c.tables[tableKey]
	columns := tableInfo.Columns
	columnsLen := len(columns)
	event := new(Event)
	event.Table = table
	event.DataBase = database

	// 遍历表数据
	rows := rowsEvent.Rows
	if len(rows) >= 1 {
		event.Before = make(map[string]interface{}, columnsLen)
		row1 := rows[0]
		if len(row1) != columnsLen {
			return nil, errors.New("数据与字段不一致")
		}
		for i, v := range row1 {
			event.Before[columns[i].Name] = v
		}
	}
	if len(rows) >= 2 {
		event.After = make(map[string]interface{}, columnsLen)
		row1 := rows[1]
		if len(row1) != columnsLen {
			return nil, errors.New("数据与字段不一致")
		}
		for i, v := range row1 {
			event.After[columns[i].Name] = v
		}
	}
	return event, nil
}

func (c *mysqlConverter) key(database string, table string) string {
	return fmt.Sprintf("%s.%s", database, table)
}