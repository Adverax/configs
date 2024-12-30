package configs

func LoadYamlConfigs(config interface{}, files ...string) error {
	return LoadConfig(
		config,
		WithDataSources(NewYamlDataSources(files...)...),
		WithConverter(NewYamlConverter()),
	)
}

func NewYamlDataSources(files ...string) []DataSource {
	var ds []DataSource
	for i, f := range files {
		ds = append(ds, NewYamlDataSource(NewFileStream(f, WithFileStreamMustExists(i == 0))))
	}
	return ds
}
