package jsonConfig

import (
	"bytes"
	"encoding/json"
	"github.com/adverax/configs"
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

	decoder := json.NewDecoder(bytes.NewBuffer(source))
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
