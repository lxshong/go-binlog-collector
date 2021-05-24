package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
)

type Event struct {
	Table     string                 `json:"table"`
	EventType string                 `json:"event_type"`
	DataBase  string                 `json:"data_base"`
	Before    map[string]interface{} `json:"before"`
	After     map[string]interface{} `json:"after"`
}

func (e *Event) String() string {
	str, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(str)
}

// 创建转换器
func NewConverter(instance string, fromType string, config *MysqlConfig) (converter, error) {
	switch fromType {
	case MYSQL:
		cvt := new(mysqlConverter)
		cvt.instance = instance
		cvt.config = config
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
	config   *MysqlConfig
	conn     *client.Conn
}

// 转换
func (c *mysqlConverter) Convert(ev *replication.BinlogEvent) (*Event, error) {
	switch ev.Event.(type) {
	case *replication.RowsEvent:
		rowsEvent := ev.Event.(*replication.RowsEvent)
		eventType := c.getEventType(ev.Header)
		if event, err := c.rowsEventConvert(eventType, rowsEvent); err != nil {
			return nil, err
		} else {
			return event, nil
		}
	}
	return nil, nil
}

// 获取操作事件类型
func (c *mysqlConverter) getEventType(header *replication.EventHeader) string {
	switch header.EventType {
	case replication.UPDATE_ROWS_EVENTv0, replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		return UPDATE
	case replication.WRITE_ROWS_EVENTv0, replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		return INSERT
	case replication.DELETE_ROWS_EVENTv0, replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		return DELETE
	}
	return ""
}

func (c *mysqlConverter) GetTable(database string, table string) error {
	key := c.key(database, table)
	if _, ok := c.tables[key]; ok {
		return nil
	}
	config := c.config
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if c.conn == nil {
		conn, err := client.Connect(addr, config.User, config.Passwd, database)
		if err != nil {
			return err
		}
		c.conn = conn
	}
	fmt.Println("get table info of  ", database+"."+table)
	tableInfo, err := schema.NewTable(c.conn, database, table)
	if err != nil {
		return err
	}
	c.tables[key] = *tableInfo
	return nil
}

// RowsEvent 事件处理
func (c *mysqlConverter) rowsEventConvert(eventType string, rowsEvent *replication.RowsEvent) (*Event, error) {
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
	event := &Event{
		Table:     table,
		DataBase:  database,
		EventType: eventType,
	}

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
