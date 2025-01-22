package mixedConfig

import (
	"fmt"
	yamlConfig "github.com/adverax/configs/formats/yaml"
	envFetcher "github.com/adverax/core/fetchers/maps/env"
)

type MyConfigAddress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MyConfig struct {
	Address MyConfigAddress `yaml:"address"`
	Name    string          `yaml:"name"`
}

func DefaultConfig() *MyConfig {
	return &MyConfig{
		Address: MyConfigAddress{
			Host: "unknown",
			Port: 80,
		},
		Name: "unknown",
	}
}

func Example() {
	// This example demonstrates how to use Mixed loader.
	//
	// First, create loader:
	loader, err := yamlConfig.NewFileLoaderBuilder().
		WithFile("config.global.json", false).
		WithFile("config.local.json", false).
		WithSource(
			envFetcher.New(
				envFetcher.NewPrefixGuard("MYAPP_"),
				envFetcher.NewKeyPathAccumulator("_"),
			),
		).
		Build()
	if err != nil {
		panic(err)
	}

	// Then load configuration:
	config := DefaultConfig()
	err = loader.Load(config)
	if err != nil {
		panic(err)
	}

	// Now you can use config.
	// For example, print it:
	fmt.Println(*config)

	// Output:
	// {{google.com 90} My App}
}
