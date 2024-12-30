package configs

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

type YamlDataSource struct {
	stream Stream
}

func NewYamlDataSource(stream Stream) DataSource {
	return &YamlDataSource{
		stream: stream,
	}
}

func (that *YamlDataSource) Fetch() (map[string]interface{}, error) {
	data := map[string]interface{}{}

	source, err := that.stream.Fetch()
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
