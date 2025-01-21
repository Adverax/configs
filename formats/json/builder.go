package jsonConfig

import (
	"github.com/adverax/configs"
	"github.com/adverax/core/sources/json"
)

func NewFileLoaderBuilder() *configs.FileLoaderBuilder {
	return configs.NewFileLoaderBuilder().
		WithSourceBuilder(
			func(fetcher configs.Fetcher) configs.Source {
				return jsonSource.New(fetcher)
			},
		).
		WithConverter(NewConverter())
}
