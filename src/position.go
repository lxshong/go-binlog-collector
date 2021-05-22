package src

import (
	"encoding/json"
	"fmt"
	"go-binlog-collector/utils"
)

// 位置信息
type position struct {
	Instance string `json:"instance"`
	Name     string `json:"name"`
	Pos      int    `json:"pos"`
}

func NewPosition(instance string) *position {
	return &position{
		Instance: instance,
	}
}

// 从缓存中初始化
func (p *position) Init() error {

	content, err := utils.FileGetContent(p.getFileName())
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return nil
	}
	if err := json.Unmarshal([]byte(content), p); err != nil {
		return err
	}
	return nil
}

// 缓存到文件
func (p *position) Save() error {
	bs, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return utils.FilePutContent(p.getFileName(), string(bs))
}

func (p *position) getFileName() string {
	return fmt.Sprintf("./tmp/%s.pos", p.Instance)
}
