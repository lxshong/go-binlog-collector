package main

import (
	"fmt"
	"go-binlog-collector/src"
)

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
	instance := "test.tb_jd_district"
	src.AddInstanceConfig(instance, &src.InstanceConfig{
		FromType: "mysql",
		FromConfig: &src.MysqlConfig{
			Host:         "127.0.0.1",
			Port:         23306,
			User:         "root",
			Passwd:       "123456",
			DataBase:     "test",
			Table:        "tb_jd_district",
			UniqueColumn: "id",
		},
		ToType: "redis",
		ToConfig: &src.RedisConfig{
			Host:   "127.0.0.1",
			Port:   26379,
			Passwd: "",
			DB:     0,
		},
	})

	if err := src.Run(instance); err != nil {
		fmt.Println(err)
	}
}
