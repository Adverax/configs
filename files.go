package configs

import (
	"fmt"
	fileFetchers "github.com/adverax/core/fetchers/files"
)

type SourceBuilder func(Fetcher) Source

type FileLoaderBuilder struct {
	files     []string
	builder   SourceBuilder
	validator Validator
	converter Converter
}

func NewFileLoaderBuilder() *FileLoaderBuilder {
	return &FileLoaderBuilder{}
}

func (that *FileLoaderBuilder) WithSourceBuilder(builder SourceBuilder) *FileLoaderBuilder {
	that.builder = builder
	return that
}

func (that *FileLoaderBuilder) WithFile(files ...string) *FileLoaderBuilder {
	that.files = append(that.files, files...)
	return that
}

func (that *FileLoaderBuilder) WithConverter(converter Converter) *FileLoaderBuilder {
	that.converter = converter
	return that
}

func (that *FileLoaderBuilder) WithValidator(validator Validator) *FileLoaderBuilder {
	that.validator = validator
	return that
}

func (that *FileLoaderBuilder) Build() (*Loader, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	sources, err := that.newFileSources(that.builder, that.files...)
	if err != nil {
		return nil, err
	}

	return NewBuilder().
		WithSources(sources...).
		WithConverter(that.converter).
		WithValidator(that.validator).
		Build()
}

func (that *FileLoaderBuilder) checkRequiredFields() error {
	if len(that.files) == 0 {
		return ErrFieldFilesIsRequired
	}

	if that.builder == nil {
		return ErrFieldBuilderIsRequired
	}

	if that.converter == nil {
		return ErrFieldConverterIsRequired
	}

	return nil
}

func (that *FileLoaderBuilder) newFileSources(
	builder func(Fetcher) Source,
	files ...string,
) ([]Source, error) {
	var ds []Source
	for i, f := range files {
		if f == "" {
			continue
		}
		fetcher, err := fileFetchers.NewBuilder().
			WithFilename(f).
			WithMustExists(i == 0).
			Build()
		if err != nil {
			return nil, err
		}
		ds = append(ds, builder(fetcher))
	}
	return ds, nil
}

var (
	ErrFieldFilesIsRequired     = fmt.Errorf("Files are required")
	ErrFieldBuilderIsRequired   = fmt.Errorf("Builder is required")
	ErrFieldConverterIsRequired = fmt.Errorf("Converter is required")
)
