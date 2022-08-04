package envconf

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	Debug bool   `env:"DEBUG" default:"true"`
	Host  string `env:"HOST"`
	Port  int
}

type ConfigWithSecret struct {
	SecretKey string `required:"true"`
}

func TestLoader_Load__Success(t *testing.T) {
	c := struct {
		Debug bool   `env:"DEBUG" default:"true"`
		Host  string `env:"HOST"`
		Port  int
	}{}

	os.Clearenv()
	os.Setenv("HOST", "localhost")

	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "localhost", c.Host)
}

func TestLoader_Load_InvalidObjectType(t *testing.T) {
	tests := []struct {
		name   string
		object interface{}
	}{
		{name: "String", object: "string"},
		{name: "Integer", object: 1337},
		{name: "Struct", object: Config{}},
		{name: "Slice", object: []Config{{}, {}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewLoader().Load(tt.object)

			assert.ErrorIs(t, ErrInvalidObjectType, err)
		})
	}
}

func TestLoader_Load__WithPrefix(t *testing.T) {
	var c Config

	os.Clearenv()
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("APP_HOST", "192.168.0.1")

	loader := NewLoader()
	loader.SetPrefix("app")

	err := loader.Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "192.168.0.1", c.Host)
}

func TestLoader_Load__DefaultValue(t *testing.T) {
	var c Config

	os.Clearenv()

	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.True(t, c.Debug)
}

func TestLoader_Load__NoEnvTag(t *testing.T) {
	var c Config

	os.Clearenv()
	os.Setenv("PORT", "1337")

	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, 1337, c.Port)
}

func TestLoader_Load__RequiredFieldIsMissing(t *testing.T) {
	var cws ConfigWithSecret

	os.Clearenv()

	err := NewLoader().Load(&cws)

	assert.Error(t, err)
}

func TestLoader_Load__MissingValue(t *testing.T) {
	var c Config

	os.Clearenv()

	err := NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "", c.Host)
}

func TestGetEnvValue_MissingValue(t *testing.T) {
	os.Clearenv()

	fi := fieldInfo{
		Name: "Field",
		Env:  "FIELD",
	}
	v, err := getEnvValue("", fi)

	assert.Nil(t, err)
	assert.Equal(t, "", v)
}

func TestGetEnvValue_MissingRequiredField(t *testing.T) {
	os.Clearenv()

	fi := fieldInfo{
		Name:     "Field",
		Env:      "FIELD",
		Required: true,
	}
	_, err := getEnvValue("", fi)

	assert.Error(t, err)
}

func TestGetEnvValue_MissingWithDefault(t *testing.T) {
	os.Clearenv()

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
		name          string
		fieldName     string
		value         string
		expectedValue interface{}
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
	os.Clearenv()

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
	os.Clearenv()
	os.Setenv("HOSTS", "192.168.0.1,192.168.0.2")

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
	os.Clearenv()
	os.Setenv("DEBUG", "invalid")

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
	os.Clearenv()
	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "1337")
	os.Setenv("DEBUG", "true")
	os.Setenv("PROCS", "3")

	c := struct {
		Host  string
		Port  int
		Debug bool
		Procs int8
	}{}
	l := NewLoader()
	tests := []struct {
		name     string
		fi       fieldInfo
		expected interface{}
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

func getFieldValue(f reflect.Value) interface{} {
	switch f.Kind() {
	case reflect.String:
		return f.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int()

	case reflect.Bool:
		return f.Bool()
	}

	return nil
}
