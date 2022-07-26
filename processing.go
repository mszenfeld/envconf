package envconf

import (
	"errors"
	"reflect"
)

type fieldInfo struct {
	Name string
	Env string
	Default interface{}
	Required bool
}

var ErrInvalidObjectType = errors.New("invalid object type")

func Process(obj interface{}) ([]fieldInfo, error) {
	v := reflect.ValueOf(obj)

	if err := validateObjType(v); err != nil {
		return []fieldInfo{}, err
	}

	return []fieldInfo{}, nil
}

func validateObjType(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return ErrInvalidObjectType
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrInvalidObjectType
	}

	return nil
}
