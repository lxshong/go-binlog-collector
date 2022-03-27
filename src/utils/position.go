package utils

import (
	"encoding/json"
	"fmt"
)

// 位置信息
type Position struct {
	Instance string `json:"instance"`
	Name     string `json:"name"`
	Pos      int    `json:"pos"`
}

func NewPosition(instance string) *Position {
	return &Position{
		Instance: instance,
	}
}

// 从缓存中初始化
func (p *Position) Init() error {

	content, err := FileGetContent(p.getFileName())
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
func (p *Position) Save() error {
	bs, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return FilePutContent(p.getFileName(), string(bs))
}

func (p *Position) getFileName() string {
	return fmt.Sprintf("./tmp/%s.pos", p.Instance)
}
