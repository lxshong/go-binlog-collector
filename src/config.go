package src

import "errors"

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Passwd   string
	DataBase string
	Table    string
}

type RedisConfig struct {
	Host     string
	Port     int
	Passwd   string
	DB int
}

type InstanceConfig struct {
	FromType   string
	ToType     string
	FromConfig interface{}
	ToConfig   interface{}
}

var instances map[string]*InstanceConfig

func init() {
	instances = make(map[string]*InstanceConfig)
}

func getInstanceConfig(instance string) (*InstanceConfig, error) {
	if config, ok := instances[instance]; ok {
		return config, nil
	}
	return nil, errors.New("instance config unexists")
}

// 添加配置
func AddInstanceConfig(instance string, config *InstanceConfig) error {
	instances[instance] = config
	return nil
}
