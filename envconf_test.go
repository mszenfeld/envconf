package envconf_test

import (
	"testing"

	"github.com/mszenfeld/envconf"
	"github.com/mszenfeld/envconf/processing"
	"github.com/stretchr/testify/assert"
)

func TestLoader_Load__Success(t *testing.T) {
	t.Setenv("HOST", "localhost")

	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := envconf.NewLoader().Load(&c)

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
			err := envconf.NewLoader().Load(tt.object)

			assert.ErrorIs(t, processing.ErrInvalidObjectType, err)
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
	loader := envconf.NewLoader()
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
	err := envconf.NewLoader().Load(&c)

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
	err := envconf.NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, 1337, c.Port)
}

func TestLoader_Load__RequiredFieldIsMissing(t *testing.T) {
	c := struct {
		SecretKey string `required:"true"`
	}{}
	err := envconf.NewLoader().Load(&c)

	assert.Error(t, err)
}

func TestLoader_Load__MissingValue(t *testing.T) {
	c := struct {
		Host  string `env:"HOST"`
		Port  int
		Debug bool `env:"DEBUG" default:"true"`
	}{}
	err := envconf.NewLoader().Load(&c)

	assert.Nil(t, err)
	assert.Equal(t, "", c.Host)
}
