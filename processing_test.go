package envconf

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Debug     bool   `env:"DEBUG" default:"true"`
	Host      string `env:"HOST"`
	Port      int
	SecretKey string `required:"true"`
}

func TestProcess(t *testing.T) {
	info, err := Process(&TestConfig{})

	assert.Nil(t, err)
	assert.Len(t, info, 4)

	assert.ObjectsAreEqual(fieldInfo{
		Name:       "Debug",
		Env:        "DEBUG",
		Default:    "true",
		HasDefault: true,
		Required:   false,
	}, info[0])
	assert.ObjectsAreEqual(fieldInfo{
		Name:       "Host",
		Env:        "HOST",
		Default:    "",
		HasDefault: false,
		Required:   false,
	}, info[1])
	assert.ObjectsAreEqual(fieldInfo{
		Name:       "Port",
		Env:        "PORT",
		Default:    "",
		HasDefault: false,
		Required:   false,
	}, info[2])
	assert.ObjectsAreEqual(fieldInfo{
		Name:       "SecretKey",
		Env:        "SECRET_KEY",
		Default:    "",
		HasDefault: false,
		Required:   true,
	}, info[3])
}

func TestProcess_DifferentTypes(t *testing.T) {
	tests := []struct {
		name       string
		obj        interface{}
		shouldFail bool
	}{
		{name: "String", obj: "object", shouldFail: true},
		{name: "Integer", obj: 10, shouldFail: true},
		{name: "Boolean", obj: true, shouldFail: true},
		{name: "Struct", obj: TestConfig{}, shouldFail: true},
		{name: "Pointer", obj: &TestConfig{}, shouldFail: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := Process(tt.obj)

			if tt.shouldFail {
				assert.ErrorIs(t, ErrInvalidObjectType, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Empty(t, info)
		})
	}
}

func TestGetEnv(t *testing.T) {
	conf := struct {
		Debug     bool
		SecretKey string
		Host      string `env:"HOST"`
		Port      int    `env:"APP_PORT"`
	}{}
	tests := []struct {
		name     string
		fieldIdx int
		expected string
	}{
		{name: "No tag & Simple", fieldIdx: 0, expected: "DEBUG"},
		{name: "No tag & CamelCase", fieldIdx: 1, expected: "SECRET_KEY"},
		{name: "With tag & Simple", fieldIdx: 2, expected: "HOST"},
		{name: "With tag & Custom env", fieldIdx: 3, expected: "APP_PORT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := reflect.ValueOf(&conf).Elem().Type()
			fType := typ.Field(tt.fieldIdx)

			assert.Equal(t, tt.expected, getEnv(fType))
		})
	}
}

func TestGetDefault(t *testing.T) {
	conf := struct {
		Debug     bool `default:"true"`
		SecretKey string
		Host      string `default:"localhost"`
		AppName   string `default:""`
	}{}
	tests := []struct {
		name       string
		fieldIdx   int
		expected   string
		hasDefault bool
	}{
		{name: "Bool default", fieldIdx: 0, expected: "true", hasDefault: true},
		{name: "No default", fieldIdx: 1, expected: "", hasDefault: false},
		{name: "String default", fieldIdx: 2, expected: "localhost", hasDefault: true},
		{name: "Empty default", fieldIdx: 3, expected: "", hasDefault: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := reflect.ValueOf(&conf).Elem().Type()
			fType := typ.Field(tt.fieldIdx)

			def, hasDef := fType.Tag.Lookup("default")

			assert.Equal(t, tt.expected, def)
			assert.Equal(t, tt.hasDefault, hasDef)
		})
	}
}

func TestIsRequired(t *testing.T) {
	conf := struct {
		Debug     bool
		SecretKey string `required:"true"`
		Host      string `required:"false"`
		Port      int    `required:"invalid"`
	}{}
	tests := []struct {
		name     string
		fieldIdx int
		expected bool
	}{
		{name: "Implicit", fieldIdx: 0, expected: false},
		{name: "Explicit True", fieldIdx: 1, expected: true},
		{name: "Explicit False", fieldIdx: 2, expected: false},
		{name: "Explicit Invalid", fieldIdx: 3, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := reflect.ValueOf(&conf).Elem().Type()
			fType := typ.Field(tt.fieldIdx)

			assert.Equal(t, tt.expected, isRequired(fType))
		})
	}
}
