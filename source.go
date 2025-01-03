package configs

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

type Source interface {
	Fetch() (map[string]interface{}, error)
}

type YamlSource struct {
	fetcher Fetcher
}

func NewYamlSource(fetcher Fetcher) Source {
	return &YamlSource{
		fetcher: fetcher,
	}
}

func (that *YamlSource) Fetch() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	source, err := that.fetcher.Fetch()
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(bytes.NewBuffer(source))
	err = decoder.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
