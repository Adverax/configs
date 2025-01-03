package configs

import (
	"io"
	"os"
)

type Fetcher interface {
	Fetch() ([]byte, error)
}

type FileFetcherOptions struct {
	mustExists bool
}

func WithFileMustExists(mustExist bool) func(opts *FileFetcherOptions) {
	return func(opts *FileFetcherOptions) {
		opts.mustExists = mustExist
	}
}

type FileFetcherOption func(*FileFetcherOptions)

type FileFetcher struct {
	filename string
	options  FileFetcherOptions
}

func NewFileFetcher(filename string, options ...FileFetcherOption) Fetcher {
	var opts FileFetcherOptions
	for _, opt := range options {
		opt(&opts)
	}

	return &FileFetcher{
		filename: filename,
		options:  opts,
	}
}

func (that *FileFetcher) Fetch() ([]byte, error) {
	file, err := os.Open(that.filename)
	if err != nil {
		if os.IsNotExist(err) && !that.options.mustExists {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

type MemoryFetcher struct {
	data []byte
}

func NewMemoryFetcher(data []byte) Fetcher {
	return &MemoryFetcher{
		data: data,
	}
}

func (that *MemoryFetcher) Fetch() ([]byte, error) {
	return that.data, nil
}
