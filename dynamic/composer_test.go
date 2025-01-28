package dynConfigs

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	type TestStruct struct {
		BaseConfig
		BoolField     Boolean
		IntField      Integer
		FloatField    Float
		StringField   String
		StringsField  Strings
		DurationField Duration
	}

	obj := &TestStruct{}
	factory := NewFieldFactory()
	composer := NewComposer(factory)

	composer.Init(obj)

	if obj.BoolField == nil {
		t.Error("BoolField not initialized")
	}
	if obj.IntField == nil {
		t.Error("IntField not initialized")
	}
	if obj.FloatField == nil {
		t.Error("FloatField not initialized")
	}
	if obj.StringField == nil {
		t.Error("StringField not initialized")
	}
	if obj.StringsField == nil {
		t.Error("StringsField not initialized")
	}
	if obj.DurationField == nil {
		t.Error("DurationField not initialized")
	}
}

func TestAssign(t *testing.T) {
	type NestedStruct struct {
		IntField Integer
	}

	type TestStruct struct {
		BaseConfig
		BoolField     Boolean
		IntField      Integer
		FloatField    Float
		StringField   String
		StringsField  Strings
		DurationField Duration
		Nested        NestedStruct
	}

	src := &TestStruct{
		BoolField:     NewBoolean(true),
		IntField:      NewInteger(42),
		FloatField:    NewFloat(3.14),
		StringField:   NewString("test"),
		StringsField:  NewStrings([]string{"a", "b"}),
		DurationField: NewDuration(time.Second),
		Nested:        NestedStruct{IntField: NewInteger(55)},
	}

	dst := &TestStruct{}
	factory := NewFieldFactory()
	composer := NewComposer(factory)
	composer.Init(src)
	composer.Init(dst)
	composer.Assign(dst, src)

	if v, _ := dst.BoolField.Get(context.Background()); v != true {
		t.Error("BoolField not assigned correctly")
	}
	if v, _ := dst.IntField.Get(context.Background()); v != 42 {
		t.Error("IntField not assigned correctly")
	}
	if v, _ := dst.FloatField.Get(context.Background()); v != 3.14 {
		t.Error("FloatField not assigned correctly")
	}
	if v, _ := dst.StringField.Get(context.Background()); v != "test" {
		t.Error("StringField not assigned correctly")
	}
	if v, _ := dst.StringsField.Get(context.Background()); !reflect.DeepEqual(v, []string{"a", "b"}) {
		t.Error("StringsField not assigned correctly")
	}
	if v, _ := dst.DurationField.Get(context.Background()); v != time.Second {
		t.Error("DurationField not assigned correctly")
	}
	if v, _ := dst.Nested.IntField.Get(context.Background()); v != 55 {
		t.Error("Nested.IntField not assigned correctly")
	}
}
