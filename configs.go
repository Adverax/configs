package configs

import (
	"fmt"
	"reflect"
)

type DataSource interface {
	Fetch() (map[string]interface{}, error)
}

type LoadConfigOptions struct {
	sources   []DataSource
	converter Converter
}

type LoadConfigOption func(*LoadConfigOptions)

func WithDataSources(sources ...DataSource) LoadConfigOption {
	return func(opts *LoadConfigOptions) {
		opts.sources = append(opts.sources, sources...)
	}
}

func WithConverter(converter Converter) LoadConfigOption {
	return func(opts *LoadConfigOptions) {
		opts.converter = converter
	}
}

func LoadConfig(config interface{}, options ...LoadConfigOption) error {
	opts := LoadConfigOptions{
		converter: defaultConverter,
	}
	for _, opt := range options {
		opt(&opts)
	}

	var data map[string]interface{}
	for _, source := range opts.sources {
		if data == nil {
			data = make(map[string]interface{})
		}

		d, err := source.Fetch()
		if err != nil {
			return fmt.Errorf("error in source: %w", err)
		}

		override(data, d)
	}

	return opts.converter.Convert(data, config)
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
