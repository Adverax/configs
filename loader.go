package configs

import (
	"fmt"
)

type Fetcher interface {
	Fetch() ([]byte, error)
}

type Source interface {
	Fetch() (map[string]interface{}, error)
}

type Converter interface {
	Convert(src, dst interface{}) error
}

type Validator interface {
	Validate(config interface{}) error
}

type Loader struct {
	sources   []Source
	converter Converter
	validator Validator
}

func (that *Loader) Load(config interface{}) error {
	data := make(map[string]interface{})
	for _, source := range that.sources {
		d, err := source.Fetch()
		if err != nil {
			return fmt.Errorf("error in source: %w", err)
		}

		override(data, d)
	}

	err := that.converter.Convert(data, config)
	if err != nil {
		return fmt.Errorf("error convert config: %w", err)
	}

	if that.validator != nil {
		err = that.validator.Validate(config)
		if err != nil {
			return fmt.Errorf("error validate config: %w", err)
		}
	}

	return nil
}
