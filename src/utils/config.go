package utils

import (
	"gopkg.in/yaml.v2"
)

const (
	MYSQL = "mysql"
	REDIS = "redis"
)

const (
	UPDATE = "update"
	INSERT = "insert"
	DELETE = "delete"
)

type MysqlConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
}

type RedisConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Passwd string `yaml:"passwd"`
	DB     int    `yaml:"db"`
}

type InstanceConfig struct {
	Instance string                    `yaml:"instance"`
	FromType string                    `yaml:"from"`
	ToType   string                    `yaml:"to"`
	Mysql    *MysqlConfig              `yaml:"mysql"`
	Redis    *RedisConfig              `yaml:"redis"`
	Rules    map[string]map[string]int `yaml:"rules"`
}

// 解析配置文件
func ParseConfig(file string) (*InstanceConfig, error) {
	content, err := FileGetContent(file)
	if err != nil {
		return nil, err
	}
	config := &InstanceConfig{}
	if err := yaml.Unmarshal([]byte(content), config); err != nil {
		return nil, err
	}
	return config, nil
}
