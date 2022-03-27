package utils

import "encoding/json"

type Event struct {
	Table     string                 `json:"table"`
	EventType string                 `json:"event_type"`
	DataBase  string                 `json:"data_base"`
	Before    map[string]interface{} `json:"before"`
	After     map[string]interface{} `json:"after"`
}

func (e *Event) String() string {
	str, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(str)
}
