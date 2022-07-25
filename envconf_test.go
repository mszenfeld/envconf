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
