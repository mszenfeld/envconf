package processing_test

import (
	"testing"

	"github.com/mszenfeld/envconf/processing"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Host      string `env:"HOST"`
	SecretKey string `required:"true"`
	Port      int
	Debug     bool `env:"DEBUG" default:"true"`
}

func TestProcess(t *testing.T) {
	info, err := processing.Process(&TestConfig{})

	assert.Nil(t, err)
	assert.Len(t, info, 4)

	assert.ObjectsAreEqual(processing.FieldInfo{
		Name:       "Debug",
		Env:        "DEBUG",
		Default:    "true",
		HasDefault: true,
		Required:   false,
	}, info[0])
	assert.ObjectsAreEqual(processing.FieldInfo{
		Name:       "Host",
		Env:        "HOST",
		Default:    "",
		HasDefault: false,
		Required:   false,
	}, info[1])
	assert.ObjectsAreEqual(processing.FieldInfo{
		Name:       "Port",
		Env:        "PORT",
		Default:    "",
		HasDefault: false,
		Required:   false,
	}, info[2])
	assert.ObjectsAreEqual(processing.FieldInfo{
		Name:       "SecretKey",
		Env:        "SECRET_KEY",
		Default:    "",
		HasDefault: false,
		Required:   true,
	}, info[3])
}

func TestProcess_DifferentTypes(t *testing.T) {
	tests := []struct {
		obj        interface{}
		name       string
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
			info, err := processing.Process(tt.obj)

			if tt.shouldFail {
				assert.ErrorIs(t, processing.ErrInvalidObjectType, err)
				assert.Empty(t, info)
			} else {
				assert.Nil(t, err)
				assert.NotEmpty(t, info)
			}
		})
	}
}
