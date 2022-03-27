package dispense

import (
	"go-binlog-collector/src/utils"
)

// 创建RedisDispenser
func newStdioDispenser(instance string, config *utils.RedisConfig) (dispenser, error) {
	// 创建dispenser
	return &stdioDispenser{
		prekey: instance,
	}, nil
}

type stdioDispenser struct {
	prekey string
}

// 发送
func (d *stdioDispenser) Do(event *utils.Event) error {
	utils.Logger.Println(event)
	return nil
}
