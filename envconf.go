package envconf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mszenfeld/envconf/processing"
)

var ErrUnsupportedType = errors.New("unsupported field type")

type Loader struct {
	prefix string
}

// NewLoader creates a new instance of the `Loader`.
func NewLoader() *Loader {
	return &Loader{}
}

// Load loads values from the environment variables and pass them to the
// provided object.
func (l *Loader) Load(obj interface{}) error {
	fieldInfos, err := processing.Process(obj)
	if err != nil {
		return err
	}

	for i := 0; i < len(fieldInfos); i++ {
		if err = l.loadField(obj, fieldInfos[i]); err != nil {
			return err
		}
	}

	return nil
}

// SetPrefix sets a new value for `prefix`.
func (l *Loader) SetPrefix(p string) {
	l.prefix = strings.ToUpper(p)
}

// Prefix returns current `prefix` value.
func (l *Loader) Prefix() string {
	return l.prefix
}

// loadField gets a value of the environment variable and sets it to the field.
func (l *Loader) loadField(obj interface{}, fi processing.FieldInfo) error {
	f := reflect.ValueOf(obj).Elem().FieldByName(fi.Name)

	v, err := getEnvValue(l.prefix, fi)
	if err != nil {
		return err
	}

	if err = setFieldValue(f, v); err != nil {
		return err
	}

	return nil
}

// getEnvValue looks for proper environment variable and gets it value.
//
// If environment variable is not available, but default value was provided,
// function will return it. In other cases it will return empty string.
//
// If value is required but environment variable does not exist, function will
// return an error.
func getEnvValue(prefix string, fi processing.FieldInfo) (string, error) {
	envName := fi.Env
	if len(prefix) > 0 {
		envName = fmt.Sprintf("%s_%s", prefix, envName)
	}
	v, ok := os.LookupEnv(envName)

	if !ok {
		if fi.Required && !fi.HasDefault {
			return "", fmt.Errorf("missing value for the required field: %s", fi.Name)
		}
		if fi.HasDefault {
			v = fi.Default
		}
	}

	return v, nil
}

// setFieldValue sets value to the provided field.
//
// It's detecting kind of the provided field and cast received string according
// to it. If it is not possible to cast string `s` to the proper value, function
// will return an error.
func setFieldValue(f reflect.Value, s string) error {
	switch f.Kind() { //nolint:exhaustive // There is no need to include all missing reflect cases
	case reflect.String:
		f.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// If string is empty use zero
		if len(s) == 0 {
			s = "0"
		}

		v, err := strconv.ParseInt(s, 0, f.Type().Bits())
		if err != nil {
			return err
		}
		f.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// If string is empty use zero
		if len(s) == 0 {
			s = "0"
		}

		v, err := strconv.ParseUint(s, 0, f.Type().Bits())
		if err != nil {
			return err
		}
		f.SetUint(v)

	case reflect.Float32, reflect.Float64:
		// If string is empty use zero
		if len(s) == 0 {
			s = "0"
		}

		v, err := strconv.ParseFloat(s, f.Type().Bits())
		if err != nil {
			return err
		}
		f.SetFloat(v)

	case reflect.Bool:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		f.SetBool(v)

	default:
		return ErrUnsupportedType
	}

	return nil
}
