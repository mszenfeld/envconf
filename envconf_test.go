package envconf

import (
	"os"
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
	var c Config

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
		{name: "Slice", object: []Config{Config{}, Config{}}},
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

func TestLoader_loadField__MissingValue(t *testing.T) {
	os.Clearenv()

	fi := fieldInfo{
		Name: "Field",
		Env:  "FIELD",
	}
	v, err := getEnvValue("", fi)

	assert.Nil(t, err)
	assert.Equal(t, "", v)
}

func TestLoader_loadField(t *testing.T) {

}
