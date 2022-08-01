package envconf

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Loader struct {
	prefix string
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(obj interface{}) error {
	fieldInfos, err := Process(obj)
	if err != nil {
		return err
	}

	for i := 0; i < len(fieldInfos); i++ {
		l.loadField(obj, fieldInfos[i])
	}

	return nil
}

func (l *Loader) SetPrefix(p string) {
	l.prefix = strings.ToUpper(p)
}

func (l *Loader) Prefix() string {
	return l.prefix
}

func (l *Loader) loadField(obj interface{}, fi fieldInfo) error {
	reflect.ValueOf(obj).Elem().FieldByName(fi.Name)
	_, err := getEnvValue(l.prefix, fi)
	if err != nil {
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
func getEnvValue(prefix string, fi fieldInfo) (string, error) {
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

func setFieldValue(f reflect.Value, v interface{}) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(v.(string))

	case reflect.Int:
		f.SetInt(v.(int64))

	case reflect.Bool:
		f.SetBool(v.(bool))
	}

	return nil
}
