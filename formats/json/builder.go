package jsonConfig

import (
	"github.com/adverax/configs"
	"github.com/adverax/core/fetchers/maps/json"
)

func NewFileLoaderBuilder() *configs.FileLoaderBuilder {
	return configs.NewFileLoaderBuilder().
		WithSourceBuilder(
			func(fetcher configs.Fetcher) configs.Source {
				return jsonFetcher.New(fetcher)
			},
		).
		WithConverter(NewConverter())
}
