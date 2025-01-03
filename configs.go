package configs

import (
	"fmt"
	"reflect"
)

type Loader struct {
	sources   []Source
	converter Converter
	validator Validator
}

func NewLoader() *Loader {
	return &Loader{}
}

func (that *Loader) WithSources(sources ...Source) *Loader {
	that.sources = append(that.sources, sources...)
	return that
}

func (that *Loader) WithConverter(converter Converter) *Loader {
	that.converter = converter
	return that
}

func (that *Loader) WithValidator(validator Validator) *Loader {
	that.validator = validator
	return that
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

func override(a, b map[string]interface{}) {
	for k, v := range b {
		if av, ok := a[k]; ok {
			if reflect.TypeOf(v) == reflect.TypeOf(av) {
				switch v.(type) {
				case map[string]interface{}:
					override(av.(map[string]interface{}), v.(map[string]interface{}))
				case []interface{}:
					a[k] = v
				default:
					a[k] = v
				}
			} else {
				a[k] = v
			}
		} else {
			a[k] = v
		}
	}
}
