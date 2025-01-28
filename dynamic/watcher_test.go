package dynConfigs

import (
	"context"
	"fmt"
	"github.com/adverax/configs"
	"github.com/adverax/fetchers/maps/maps"
	"time"
)

type MyConfig struct {
	BaseConfig
	Name     String  `config:"name"`
	Interval Integer `config:"interval"`
	StartAt  Time    `config:"start_at"`
}

func DefaultConfig() *MyConfig {
	return &MyConfig{
		Name: NewString("unknown"),
	}
}

func Example() {
	// This example demonstrates how to use dynamic loading configs.
	//
	// First, create source:
	source := maps.Engine{
		"name":     "My App",
		"interval": int64(10),
		"start_at": time.Now().Format(time.RFC3339),
	}

	// Then create loader:
	loader, err := configs.NewBuilder().
		WithSource(source).
		WithDistinct(true).
		Build()
	if err != nil {
		panic(err)
	}

	// Then load initial configuration:
	config := DefaultConfig()
	composer := NewComposer(NewFieldFactory())
	composer.Init(config)
	err = loader.Load(config)
	if err != nil {
		panic(err)
	}

	// Then create watcher:
	watcher, err := NewWatcherBuilder().
		WithConfig(config).
		WithLoader(loader).
		WithInterval(10 * time.Second).
		WithNewConfig(func() Config {
			return DefaultConfig()
		}).
		Build()
	if err != nil {
		panic(err)
	}

	// Start watcher:
	watcher.Start()
	defer watcher.Close()

	// Now you can use dynamic config.
	usingConfig(context.Background(), config)

	// Output:
	// My App 10
}

func usingConfig(ctx context.Context, config *MyConfig) {
	// For example, print it:
	name, _ := config.Name.Get(ctx)
	interval, _ := config.Interval.Get(ctx)
	fmt.Println(name, interval)
}
