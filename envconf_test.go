package envconf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyConfig struct {
	Host string `env:"HOST"`
}

func TestLoader_Load__Success(t *testing.T) {
	var mc MyConfig

	os.Clearenv()
	os.Setenv("HOST", "localhost")

	err := NewLoader().Load(&mc)

	assert.Nil(t, err)
	assert.Equal(t, "localhost", mc.Host)
}

func TestLoader_Load__WithPrefix(t *testing.T) {

}
