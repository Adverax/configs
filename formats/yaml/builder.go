package yamlConfig

import (
	"github.com/adverax/configs"
	yamlSource "github.com/adverax/core/sources/yaml"
)

type Source struct {
	fetcher configs.Fetcher
}

func NewFileLoaderBuilder() *configs.FileLoaderBuilder {
	return configs.NewFileLoaderBuilder().
		WithSourceBuilder(
			func(fetcher configs.Fetcher) configs.Source {
				return yamlSource.New(fetcher)
			},
		).
		WithConverter(NewConverter())
}
