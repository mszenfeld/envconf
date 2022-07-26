package envconf

import (
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
	assert.Len(t, info, 3)

	assert.ObjectsAreEqual(fieldInfo{Name: "Debug", Env: "DEBUG", Default: true, Required: false}, info[0])
	assert.ObjectsAreEqual(fieldInfo{Name: "Host", Env: "HOST", Default: nil, Required: false}, info[1])
	assert.ObjectsAreEqual(fieldInfo{Name: "Port", Env: "PORT", Default: nil, Required: true}, info[2])
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := Process(tt.obj)

			assert.ErrorIs(t, ErrInvalidObjectType, err)
			assert.Empty(t, info)
		})
	}
}
