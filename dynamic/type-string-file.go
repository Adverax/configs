package dynConfigs

import (
	"context"
	"github.com/adverax/configs"
	"os"
	"reflect"
	"strings"
	"sync"
)

type FileString interface {
	Get(ctx context.Context) (string, error)
}

type FileStringField struct {
	config Config
	sync.RWMutex
	defValue string
	value    string
	filename string
}

func NewFileString(filename string, defValue string) *FileStringField {
	return &FileStringField{
		filename: filename,
		value:    defValue,
		defValue: defValue,
	}
}

func (that *FileStringField) Init(c Config) {
	that.config = c
}

func (that *FileStringField) Set(ctx context.Context, value string) error {
	that.config.Lock()
	defer that.config.Unlock()

	return that.Let(ctx, value)
}

func (that *FileStringField) Get(ctx context.Context) (string, error) {
	that.config.RLock()
	defer that.config.RUnlock()

	return that.Fetch(ctx)
}

func (that *FileStringField) Let(ctx context.Context, value string) error {
	that.Lock()
	defer that.Unlock()

	err := os.WriteFile(that.filename, []byte(value), 0644)
	if err != nil {
		return err
	}

	that.value = value
	return nil
}

func (that *FileStringField) Fetch(ctx context.Context) (string, error) {
	that.RLock()
	defer that.RUnlock()

	if that.value != that.defValue {
		return that.value, nil
	}

	_, err := os.Stat(that.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return that.defValue, nil
		}
		return "", err
	}

	data, err := os.ReadFile(that.filename)
	if err != nil {
		return "", err
	}

	that.value = strings.TrimSpace(string(data))
	return that.value, nil
}

func (that *FileStringField) String() string {
	val, _ := that.Fetch(context.Background())
	return val
}

type FileStringTypeHandler struct {
	configs.StringTypeHandler
}

func (that *FileStringTypeHandler) Get(ctx context.Context, field interface{}) (interface{}, error) {
	if f, ok := field.(FileStringField); ok {
		return f.Get(ctx)
	}

	return nil, nil
}

func (that *FileStringTypeHandler) New(conf Config) interface{} {
	panic("implement me")
}

func init() {
	configs.RegisterHandler(reflect.TypeOf((*FileString)(nil)).Elem(), &FileStringTypeHandler{})
}
