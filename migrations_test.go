package configs

import (
	"fmt"
	"github.com/adverax/fetchers/maps/maps"
	"time"
)

type MyConfigAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type MyConfig struct {
	Address  MyConfigAddress `json:"address"`
	Name     string          `json:"name"`
	Interval time.Duration   `json:"interval"`
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
	// This example demonstrates how to use loader with migrations.
	//
	// First, create source:
	source := maps.Engine{
		"address": map[string]interface{}{
			"host": "google.com",
			"port": 90,
		},
		"name":     "My App",
		"interval": 10,
	}

	// Then create migrator:
	migrator := NewMigrator()
	migrator.Add(
		"1",
		func(data map[string]interface{}) error {
			if v, ok := data["interval"]; ok {
				data["interval"] = int64(time.Duration(v.(int)) * time.Second)
			}
			return nil
		},
	)

	// Then create loader:
	loader, err := NewBuilder().
		WithSource(NewSourceWithMigration(source, migrator)).
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
	// {{google.com 90} My App 10s}
}
