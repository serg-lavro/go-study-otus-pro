package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func validateInt(name string, req reflect.StructTag) {
	fmt.Println("validating INT ", name, " req ", req)
}

func validateString(name string, req reflect.StructTag) {
	fmt.Println("validating STRING ", name, " req ", req)
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		st := reflect.TypeOf(v)
		for i := 0; i < st.NumField(); i++ {
			// check if tag 'validate'
			f := st.Field(i)
			fv := rv.Field(i)
			ft := f.Type
			switch ft.Kind() {
			case reflect.Int:
				fmt.Println("ININTTT: ", fv.Int())
				validateInt(f.Name, f.Tag)
			case reflect.String:
				validateString(f.Name, f.Tag)
			}
		}

	}
	return nil
}
