package yamlConfig

import (
	"github.com/adverax/configs"
)

func NewFileLoaderBuilder() *configs.FileLoaderBuilder {
	return configs.NewFileLoaderBuilder().
		WithSourceBuilder(NewSource).
		WithConverter(NewConverter())
}
