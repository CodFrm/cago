package file

import (
	"io/ioutil"

	"github.com/codfrm/cago/config"
)

type fileSource struct {
	path          string
	config        map[string]interface{}
	serialization Serialization
}

func NewSource(filename string, serialization Serialization) (config.Source, error) {
	f := &fileSource{path: filename, serialization: serialization, config: make(map[string]interface{})}
	b, err := f.Read()
	if err != nil {
		return nil, err
	}
	if err := f.serialization.Unmarshal(b, &f.config); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *fileSource) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func (f *fileSource) Scan(key string, value interface{}) error {
	var b, err = f.serialization.Marshal(f.config[key])
	if err != nil {
		return err
	}
	return f.serialization.Unmarshal(b, value)
}
