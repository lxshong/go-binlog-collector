package src

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
	FromType string              `yaml:"from"`
	ToType   string              `yaml:"to"`
	Mysql    *MysqlConfig        `yaml:"mysql"`
	Redis    *RedisConfig        `yaml:"redis"`
	Rules    map[string][]string `yaml:"rules"`
}
