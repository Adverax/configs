package configs

import "fmt"

type YamlLoader struct {
	files     []string
	validator Validator
}

func NewYamlLoader() *YamlLoader {
	return &YamlLoader{}
}

func (that *YamlLoader) Files(files ...string) *YamlLoader {
	that.files = append(that.files, files...)
	return that
}

func (that *YamlLoader) Validator(validator Validator) *YamlLoader {
	that.validator = validator
	return that
}

func (that *YamlLoader) Load(config interface{}) error {
	if len(that.files) == 0 {
		return fmt.Errorf("Primary configuration file is required")
	}

	return NewLoader().
		WithSources(that.newSources(that.files...)...).
		WithConverter(NewYamlConverter()).
		WithValidator(that.validator).
		Load(config)
}

func (that *YamlLoader) newSources(files ...string) []Source {
	var ds []Source
	for i, f := range files {
		if i != 0 && f == "" {
			continue
		}
		ds = append(ds, NewYamlSource(NewFileFetcher(f, WithFileMustExists(i == 0))))
	}
	return ds
}
