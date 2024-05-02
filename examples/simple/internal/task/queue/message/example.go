package message

import (
	"encoding/json"
	"strconv"
)

type ExampleMsg struct {
	Time int64 `json:"time"`
}

func (e *ExampleMsg) Marshal() []byte {
	return []byte("{\"time\":" + strconv.FormatInt(e.Time, 10) + "}")
}

func (e *ExampleMsg) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}
