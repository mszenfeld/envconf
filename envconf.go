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
