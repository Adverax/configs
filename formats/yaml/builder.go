package yamlConfig

import (
	"github.com/adverax/configs"
	"github.com/adverax/fetchers/maps/yaml"
)

type Source struct {
	fetcher configs.Fetcher
}

func NewFileLoaderBuilder() *configs.FileLoaderBuilder {
	return configs.NewFileLoaderBuilder().
		WithSourceBuilder(
			func(fetcher configs.Fetcher) configs.Source {
				return yamlFetcher.New(fetcher)
			},
		).
		WithConverter(NewConverter())
}
