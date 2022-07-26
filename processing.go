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

	elem := v.Elem()

	return processFields(elem, elem.Type()) 
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

func processFields(v reflect.Value, t reflect.Type) ([]fieldInfo, error) {
	fieldInfos := make([]fieldInfo, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fType := t.Field(i)

		fieldInfo := fieldInfo{
			Name: fType.Name,
		}

		fieldInfos = append(fieldInfos, fieldInfo)
	}

	return fieldInfos, nil
}
