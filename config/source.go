package config

import "io/ioutil"

type Source interface {
	Read() ([]byte, error)
}

type fileSource struct {
	path string
}

func (f *fileSource) Read() ([]byte, error) {
	return ioutil.ReadFile(f.path)
}

func YamlFile(name string) Source {
	return &fileSource{path: name}
}
