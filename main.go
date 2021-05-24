package main

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/schema"
	"go-binlog-collector/src"

	//"go-binlog-collector/src"
	"flag"
)

var instance *string

func init() {
	instance = flag.String("i","","instance of mysql")
	flag.Parse()
}

func main() {

	//addr := fmt.Sprintf("%s:%s", "127.0.0.1", "23306")
	//conn,_ := client.Connect(addr,"root","123456","test")
	//table,err := schema.NewTable(conn,"test","tb_jd_district")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for _, column := range table.Columns {
	//	fmt.Println(column.Name)
	//	fmt.Println(column.Type)
	//	fmt.Println(column.RawType)
	//	fmt.Println(column.IsUnsigned)
	//}
	//instance := "test"
	//src.AddInstanceConfig(instance, &src.InstanceConfig{
	//	FromType: "mysql",
	//	FromConfig: &src.MysqlConfig{
	//		Host:         "127.0.0.1",
	//		Port:         23306,
	//		User:         "root",
	//		Passwd:       "123456",
	//		DataBase:     "test",
	//		Table:        "tb_jd_district",
	//		UniqueColumn: "id",
	//	},
	//	ToType: "redis",
	//	ToConfig: &src.RedisConfig{
	//		Host:   "127.0.0.1",
	//		Port:   26379,
	//		Passwd: "",
	//		DB:     0,
	//	},
	//})
	//
	//if err := src.Run(instance); err != nil {
	//	fmt.Println(err)
	//}

	if err := src.Run(*instance); err != nil {
		fmt.Println(err)
	}

	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 23306)
	conn, _ := client.Connect(addr, "root", "123456", "mysql")
	tableInfo, err := schema.NewTable(conn, "mysql", "time_zone_transition_type")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tableInfo)
}
