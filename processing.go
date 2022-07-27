package envconf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/camelcase"
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

		// Skip this field if it's not possible to set it
		if !f.CanSet() {
			continue
		}

		fieldInfo := fieldInfo{
			Name: fType.Name,
			Env: getEnv(fType),
			Required: isRequired(fType),
		}

		fieldInfos = append(fieldInfos, fieldInfo)
	}

	return fieldInfos, nil
}

// getEnv returns name of the environment variable associated with the provided
// struct field.
// 
// If provided field has `env` tag, its value will be returned. In other cases
// function getting name of the environment variable from the field name.
func getEnv(fType reflect.StructField) string {
	if v := fType.Tag.Get("env"); len(v) > 0 {
		return v
	}	
	
	var wl []string
	
	for _, word := range camelcase.Split(fType.Name) {
		wl = append(wl, strings.ToUpper(word))
	}

	return strings.Join(wl, "_")
}

func isRequired(fType reflect.StructField) bool {
	v := fType.Tag.Get("required")

	if len(v) == 0 || !isBool(v) {
		return false
	} 
	isReq, _ := strconv.ParseBool(v)

	return isReq
}

func isBool(v string) bool {
	v = strings.ToLower(v)

	return v == "true" || v == "false"
}
