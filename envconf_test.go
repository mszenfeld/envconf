package envconf

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoader_Load__Success(t *testing.T) {
	t.Setenv("HOST", "localhost")

	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "localhost", c.Host)
}

func TestLoader_Load_InvalidObjectType(t *testing.T) {
	type conf struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}

	tests := []struct {
		object interface{}
		name   string
	}{
		{name: "String", object: "string"},
		{name: "Integer", object: 1337},
		{name: "Struct", object: conf{}},
		{name: "Slice", object: []conf{{}, {}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewLoader().Load(tt.object)

			assert.ErrorIs(t, ErrInvalidObjectType, err)
		})
	}
}

func TestLoader_Load__WithPrefix(t *testing.T) {
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("APP_HOST", "192.168.0.1")

	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	loader := NewLoader()
	loader.SetPrefix("app")

	err := loader.Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "192.168.0.1", c.Host)
}

func TestLoader_Load__DefaultValue(t *testing.T) {
	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.True(t, c.Debug)
}

func TestLoader_Load__NoEnvTag(t *testing.T) {
	t.Setenv("PORT", "1337")

	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, 1337, c.Port)
}

func TestLoader_Load__RequiredFieldIsMissing(t *testing.T) {
	c := struct {
		SecretKey string `required:"true"`
	}{}
	err := NewLoader().Load(&c)

	assert.Error(t, err)
}

func TestLoader_Load__MissingValue(t *testing.T) {
	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "", c.Host)
}

func TestGetEnvValue_MissingValue(t *testing.T) {
	fi := fieldInfo{
		Name: "Field",
		Env:  "FIELD",
	}
	v, err := getEnvValue("", fi)

	assert.Nil(t, err)
	assert.Equal(t, "", v)
}

func TestGetEnvValue_MissingRequiredField(t *testing.T) {
	fi := fieldInfo{
		Name:     "Field",
		Env:      "FIELD",
		Required: true,
	}
	_, err := getEnvValue("", fi)

	assert.Error(t, err)
}

func TestGetEnvValue_MissingWithDefault(t *testing.T) {
	fi := fieldInfo{
		Name:       "Field",
		Env:        "FIELD",
		HasDefault: true,
		Default:    "MyDefaultValue",
	}
	v, err := getEnvValue("", fi)

	assert.Nil(t, err)
	assert.Equal(t, "MyDefaultValue", v)
}

func TestSetFieldValue_UnsupportedType(t *testing.T) {
	c := struct {
		Hosts []string
	}{}
	f := reflect.ValueOf(&c).Elem().Field(0)

	err := setFieldValue(f, "localhost,192.160.0.1")

	assert.ErrorIs(t, ErrUnsupportedType, err)
}

func TestSetFieldValue_InvalidType(t *testing.T) {
	c := struct {
		Host  string
		Port  int
		Debug bool
	}{}
	f := reflect.ValueOf(&c).Elem().Field(2)

	err := setFieldValue(f, "10")

	assert.Error(t, err)
	assert.Equal(t, "", c.Host)
}

func TestSetFieldValue(t *testing.T) {
	c := struct {
		Host  string
		Port  int
		Debug bool
	}{}
	tests := []struct {
		expectedValue interface{}
		name          string
		fieldName     string
		value         string
	}{
		{name: "String", fieldName: "Host", value: "localhost", expectedValue: "localhost"},
		{name: "Integer", fieldName: "Port", value: "1337", expectedValue: int64(1337)},
		{name: "Boolean", fieldName: "Debug", value: "true", expectedValue: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reflect.ValueOf(&c).Elem().FieldByName(tt.fieldName)
			err := setFieldValue(f, tt.value)

			assert.Nil(t, err)
			assert.Equal(t, tt.expectedValue, getFieldValue(f))
		})
	}
}

func TestLoader_loadField__MissingValue(t *testing.T) {
	c := struct {
		Host  string
		Port  int
		Debug bool
	}{}
	l := NewLoader()

	err := l.loadField(&c, fieldInfo{
		Name: "Host",
		Env:  "HOST",
	})

	assert.Nil(t, err)
	assert.Equal(t, "", c.Host)
}

func TestLoader_loadField__UnsupportedType(t *testing.T) {
	t.Setenv("HOSTS", "192.168.0.1,192.168.0.2")

	c := struct {
		Hosts []string
		Port  int
		Debug bool
	}{}
	l := NewLoader()

	err := l.loadField(&c, fieldInfo{
		Name: "Hosts",
		Env:  "HOSTS",
	})

	assert.Error(t, err)
	assert.Empty(t, c.Hosts)
}

func TestLoader_loadField__InvalidValueType(t *testing.T) {
	t.Setenv("DEBUG", "invalid")

	c := struct {
		Host  string
		Port  int
		Debug bool
	}{}
	l := NewLoader()

	err := l.loadField(&c, fieldInfo{
		Name: "Debug",
		Env:  "DEBUG",
	})

	assert.Error(t, err)
	assert.Zero(t, c.Debug)
}

func TestLoader_loadField(t *testing.T) {
	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "1337")
	t.Setenv("DEBUG", "true")
	t.Setenv("PROCS", "3")

	c := struct {
		Host  string
		Port  int
		Debug bool
		Procs int8
	}{}
	l := NewLoader()
	tests := []struct {
		expected any
		name     string
		fi       fieldInfo
	}{
		{
			name:     "String",
			fi:       fieldInfo{Name: "Host", Env: "HOST"},
			expected: "localhost",
		},
		{
			name:     "Integer",
			fi:       fieldInfo{Name: "Port", Env: "PORT"},
			expected: int64(1337),
		},
		{
			name:     "Boolean",
			fi:       fieldInfo{Name: "Debug", Env: "DEBUG"},
			expected: true,
		},
		{
			name:     "Integer8",
			fi:       fieldInfo{Name: "Procs", Env: "PROCS"},
			expected: int64(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := reflect.ValueOf(&c).Elem().FieldByName(tt.fi.Name)

			assert.Nil(t, l.loadField(&c, tt.fi))
			assert.Equal(t, tt.expected, getFieldValue(f))
		})
	}
}

func getFieldValue(f reflect.Value) any {
	switch f.Kind() { // nolint:exhaustive // There is no need to include all missing reflect cases
	case reflect.String:
		return f.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int()

	case reflect.Bool:
		return f.Bool()

	default:
		return nil
	}
}
