package configs

import (
	"fmt"
	"sort"
)

const MigrationKey = "migration"

type Migration func(data map[string]interface{}) error

type migration struct {
	id        string
	migration Migration
}

type migrations []migration

func (that migrations) Len() int {
	return len(that)
}

func (that migrations) Less(i, j int) bool {
	return that[i].id < that[j].id
}

func (that migrations) Swap(i, j int) {
	that[i], that[j] = that[j], that[i]
}

type Migrator struct {
	migrations migrations
}

func NewMigrator() *Migrator {
	return &Migrator{}
}

func (that *Migrator) Add(id string, m Migration) {
	that.migrations = append(that.migrations, migration{id: id, migration: m})
	sort.Sort(that.migrations)
}

func (that *Migrator) Migrate(data map[string]interface{}) error {
	var migration string
	if m, ok := data[MigrationKey]; ok {
		migration = m.(string)
	}

	for _, m := range that.migrations {
		if m.id <= migration {
			continue
		}

		err := m.migration(data)
		if err != nil {
			return fmt.Errorf("error in migration %s: %w", m.id, err)
		}

		data[MigrationKey] = m.id
	}
	return nil
}

type Writer interface {
	Save(data map[string]interface{}) error
}

type SourceMigrator interface {
	Migrate(data map[string]interface{}) error
}

type SourceWithMigration struct {
	source   Source
	migrator SourceMigrator
}

func NewSourceWithMigration(source Source, migrator SourceMigrator) *SourceWithMigration {
	return &SourceWithMigration{source: source, migrator: migrator}
}

func (that *SourceWithMigration) Fetch() (map[string]interface{}, error) {
	data, err := that.source.Fetch()
	if err != nil {
		return nil, fmt.Errorf("error in source: %w", err)
	}

	err = that.migrator.Migrate(data)
	if err != nil {
		return nil, fmt.Errorf("error in migrator: %w", err)
	}

	if writer, ok := that.source.(Writer); ok {
		err := writer.Save(data)
		if err != nil {
			return nil, fmt.Errorf("error in writer: %w", err)
		}
	}

	return data, nil
}
