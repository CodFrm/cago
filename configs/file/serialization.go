package file

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type jsonSerialization struct {
}

func (j *jsonSerialization) Marshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *jsonSerialization) Unmarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

func Json() Serialization {
	return &jsonSerialization{}
}

type yamlSerialization struct {
}

func Yaml() Serialization {
	return &yamlSerialization{}
}

func (y *yamlSerialization) Marshal(i interface{}) ([]byte, error) {
	return yaml.Marshal(i)
}

func (y *yamlSerialization) Unmarshal(bytes []byte, i interface{}) error {
	return yaml.Unmarshal(bytes, i)
}
