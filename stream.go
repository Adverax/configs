package configs

import (
	"io"
	"os"
)

type Stream interface {
	Fetch() ([]byte, error)
}

type FileStreamOptions struct {
	mustExists bool
}

func WithFileStreamMustExists(mustExist bool) func(opts *FileStreamOptions) {
	return func(opts *FileStreamOptions) {
		opts.mustExists = mustExist
	}
}

type FileStreamOption func(*FileStreamOptions)

type FileStream struct {
	filename string
	options  FileStreamOptions
}

func NewFileStream(filename string, options ...FileStreamOption) Stream {
	var opts FileStreamOptions
	for _, opt := range options {
		opt(&opts)
	}

	return &FileStream{
		filename: filename,
		options:  opts,
	}
}

func (that *FileStream) Fetch() ([]byte, error) {
	file, err := os.Open(that.filename)
	if err != nil {
		if os.IsNotExist(err) && !that.options.mustExists {
			return nil, nil
		}
		panic(err)
	}
	defer file.Close()

	return io.ReadAll(file)
}
