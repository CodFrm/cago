package memory

import (
	"encoding/json"

	"github.com/codfrm/cago/configs/source"
)

type Memory struct {
	config map[string]interface{}
}

func NewSource(config map[string]interface{}) source.Source {
	if _, ok := config["debug"]; !ok {
		config["debug"] = true
	}
	return &Memory{
		config: config,
	}
}

func (e *Memory) Scan(key string, value interface{}) error {
	if v, ok := e.config[key]; ok {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, value)
	}
	return nil
}

func (e *Memory) Has(key string) (bool, error) {
	_, ok := e.config[key]
	return ok, nil
}
