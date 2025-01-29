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
	Init(obj)

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
	Init(src)
	Init(dst)
	Assign(dst, src)

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

func TestIsStaticUpdated0(t *testing.T) {
	type TestStruct struct {
		BaseConfig
		StringField String `config:"string_field,static"`
		IntField    Integer
	}

	a := &TestStruct{
		StringField: NewString("test"),
		IntField:    NewInteger(42),
	}
	Init(a)

	b := &TestStruct{
		StringField: NewString("test2"),
		IntField:    NewInteger(45),
	}
	Init(b)

	isStatic, err := isStaticUpdated(context.Background(), a, b)
	if err != nil {
		t.Error(err)
	}

	if !isStatic {
		t.Error("Static fields should be updated")
	}
}

func TestIsStaticUpdated(t *testing.T) {
	type NestedStruct struct {
		IntField Integer `config:"int_field,static"`
	}

	type TestStruct struct {
		BaseConfig
		StringField String `config:"string_field,static"`
		IntField    Integer
		Nested      NestedStruct
	}

	tests := []struct {
		name     string
		a, b     *TestStruct
		expected bool
	}{
		{
			name: "Static field updated",
			a: &TestStruct{
				StringField: NewString("test"),
				IntField:    NewInteger(42),
			},
			b: &TestStruct{
				StringField: NewString("test2"),
				IntField:    NewInteger(42),
			},
			expected: true,
		},
		{
			name: "Static nested field updated",
			a: &TestStruct{
				Nested: NestedStruct{IntField: NewInteger(55)},
			},
			b: &TestStruct{
				Nested: NestedStruct{IntField: NewInteger(56)},
			},
			expected: true,
		},
		{
			name: "No static field updated",
			a: &TestStruct{
				StringField: NewString("test"),
				IntField:    NewInteger(42),
			},
			b: &TestStruct{
				StringField: NewString("test"),
				IntField:    NewInteger(45),
			},
			expected: false,
		},
		{
			name: "No static nested field updated",
			a: &TestStruct{
				Nested: NestedStruct{IntField: NewInteger(55)},
			},
			b: &TestStruct{
				Nested: NestedStruct{IntField: NewInteger(55)},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.a)
			Init(tt.b)

			isStatic, err := isStaticUpdated(context.Background(), tt.a, tt.b)
			if err != nil {
				t.Error(err)
			}

			if isStatic != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, isStatic)
			}
		})
	}
}
