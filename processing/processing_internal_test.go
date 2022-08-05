package processing

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	conf := struct {
		SecretKey string
		Host      string `env:"HOST"`
		Port      int    `env:"APP_PORT"`
		Debug     bool
	}{}
	tests := []struct {
		name     string
		expected string
		fieldIdx int
	}{
		{name: "No tag & CamelCase", fieldIdx: 0, expected: "SECRET_KEY"},
		{name: "With tag & Simple", fieldIdx: 1, expected: "HOST"},
		{name: "With tag & Custom env", fieldIdx: 2, expected: "APP_PORT"},
		{name: "No tag & Simple", fieldIdx: 3, expected: "DEBUG"},
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
		SecretKey string
		Host      string `default:"localhost"`
		AppName   string `default:""`
		Debug     bool   `default:"true"`
	}{}
	tests := []struct {
		name       string
		expected   string
		fieldIdx   int
		hasDefault bool
	}{
		{name: "Bool default", fieldIdx: 3, expected: "true", hasDefault: true},
		{name: "No default", fieldIdx: 0, expected: "", hasDefault: false},
		{name: "String default", fieldIdx: 1, expected: "localhost", hasDefault: true},
		{name: "Empty default", fieldIdx: 2, expected: "", hasDefault: true},
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
		SecretKey string `required:"true"`
		Host      string `required:"false"`
		Port      int    `required:"invalid"`
		Debug     bool
	}{}
	tests := []struct {
		name     string
		fieldIdx int
		expected bool
	}{
		{name: "Implicit", fieldIdx: 3, expected: false},
		{name: "Explicit True", fieldIdx: 0, expected: true},
		{name: "Explicit False", fieldIdx: 1, expected: false},
		{name: "Explicit Invalid", fieldIdx: 2, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := reflect.ValueOf(&conf).Elem().Type()
			fType := typ.Field(tt.fieldIdx)

			assert.Equal(t, tt.expected, isRequired(fType))
		})
	}
}
