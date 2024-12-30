package configs

import "gopkg.in/yaml.v3"

type Converter interface {
	Convert(src, dst interface{}) error
}

type YamlConverter struct {
}

func NewYamlConverter() *YamlConverter {
	return &YamlConverter{}
}

func (t *YamlConverter) Convert(src, dst interface{}) error {
	raw, err := yaml.Marshal(src)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(raw, dst)
	if err != nil {
		return err
	}

	return nil
}

var defaultConverter = NewYamlConverter()
