package envconf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/camelcase"
)

// fieldInfo is a representation of the struct field.
type fieldInfo struct {
	Name string
	Env string
	Default string
	HasDefault bool
	Required bool
}

var ErrInvalidObjectType = errors.New("invalid object type")

// Process extracts information about fields that make up the provided object.
//
// If object is not a pointer to the struct, this function will return an error.
func Process(obj interface{}) ([]fieldInfo, error) {
	v := reflect.ValueOf(obj)

	if err := validateObjType(v); err != nil {
		return []fieldInfo{}, err
	}

	elem := v.Elem()

	return processFields(elem, elem.Type()) 
}

// validateObjType returns error for all kinds except pointer to struct.
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

// processFields extracts information about fields from the provided arguments.
func processFields(v reflect.Value, t reflect.Type) ([]fieldInfo, error) {
	fieldInfos := make([]fieldInfo, 0, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fType := t.Field(i)

		// Skip this field if it's not possible to set it
		if !f.CanSet() {
			continue
		}

		def, hasDef := fType.Tag.Lookup("default")

		fieldInfo := fieldInfo{
			Name: fType.Name,
			Env: getEnv(fType),
			Default: def,
			HasDefault: hasDef,
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

// isRequired returns bool value with information if provided field is required.
//
// By default, fields are optional. Field is required only when `required` tag 
// was explicitly provided with "true" value. In other cases function will return
// `false` even if value for `required` tag is invalid.
func isRequired(fType reflect.StructField) bool {
	v := fType.Tag.Get("required")

	if len(v) == 0 || !isBool(v) {
		return false
	} 
	isReq, _ := strconv.ParseBool(v)

	return isReq
}

// isBool returns information if given string is a proper string version of
// the boolean value.
//
// This function returns `true` only for the following values:
// - true
// - false
func isBool(v string) bool {
	v = strings.ToLower(v)

	return v == "true" || v == "false"
}
