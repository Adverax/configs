package yamlConfig

import (
	"bytes"
	"github.com/adverax/configs"
	"gopkg.in/yaml.v3"
)

type Source struct {
	fetcher configs.Fetcher
}

func NewSource(fetcher configs.Fetcher) configs.Source {
	return &Source{
		fetcher: fetcher,
	}
}

func (that *Source) Fetch() (map[string]interface{}, error) {
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
