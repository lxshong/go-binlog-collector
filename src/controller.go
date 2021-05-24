package src

import (
	"errors"
	"flag"
	"fmt"
	"go-binlog-collector/utils"
	"gopkg.in/yaml.v2"
)

var conf *string

func init() {
	conf = flag.String("conf", "./conf", "the path of config file")
}

// 执行
func Run(instance string) error {
	config, err := parseConfig(instance)
	if err != nil {
		return err
	}
	collector, err := NewCollector(instance, config)
	if err != nil {
		return err
	}
	dispenser, err := NewDispenser(instance, config)
	if err != nil {
		return err
	}
	return collector.Run(func(event *Event) error {
		return dispenser.Send(event)
	})
}

// 解析配置文件
func parseConfig(instance string) (*InstanceConfig, error) {
	file := fmt.Sprintf("%s/%s.yaml", *conf, instance)
	content, err := utils.FileGetContent(file)
	if err != nil {
		return nil, errors.New("config file do not exists")
	}
	config := &InstanceConfig{}
	if err := yaml.Unmarshal([]byte(content), config); err != nil {
		return nil, err
	}
	return config, nil
}
