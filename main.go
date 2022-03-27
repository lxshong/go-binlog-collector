package main

import (
	"context"
	"flag"
	"fmt"
	"go-binlog-collector/src"
)

var conf *string

func init() {
	conf = flag.String("conf", "./conf", "the path of config file")
}

func main() {
	if err := src.Do(context.Background(), *conf); err != nil {
		fmt.Println(err)
	}
}
