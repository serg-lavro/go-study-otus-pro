package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
	"slices"
	"strconv"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

type ValidationRequirenment struct {
	reqType string
	reqVal string
}

func parseRequirenments(tag reflect.StructTag) []ValidationRequirenment {
	validateTag := tag.Get("validate")
    if validateTag == "" {
        return nil
    }
    parts := strings.Split(validateTag, "|")

	reqs := make([]ValidationRequirenment, 0)

	for _, p := range(parts) {
		reqFields := strings.Split(p, ":")
		reqs = append(reqs, ValidationRequirenment{ reqFields[0], reqFields[1]})
	}
    return reqs
}

func validateIn(fieldVal, set string) bool {
	allowedSet := strings.Split(set, ",")
	return slices.Contains(allowedSet, fieldVal)
}

func validationFailedInt(fieldName string, val int, tag reflect.StructTag) (ValidationError, bool) {
	reqs := parseRequirenments(tag)
	fmt.Println("[VALIDATE] ", fieldName, val, reqs)
	for _, r := range(reqs) {
		switch r.reqType {
		case "min":
			if strconv.Atoi() < val {
				fmt.Println("min failed: ", val)
			}
		case "max":
			if val > strconv.Atoi() {
				fmt.Println("max failed: ", val)
			}
		case "in":
			if !validateIn(strconv.Itoa(val), r.reqVal) {
				fmt.Println("in failed: ", val)
			}
		default:
		}
	}
	return ValidationError{}, false
}

func validationFailedString(fieldName string, val string, tag reflect.StructTag) (ValidationError, bool) {
	reqs := parseRequirenments(tag)
	fmt.Println("[VALIDATE] ", fieldName, val, reqs)
	for _, r := range(reqs) {
		switch r.reqType {
		case "in":
			if !validateIn(val, r.reqVal) {
				fmt.Println("in failed: ", val)
			}
		case "ln":
		case "reqexp":
		default:
		}
	}
	return ValidationError{}, false
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		st := reflect.TypeOf(v)
		for i := 0; i < st.NumField(); i++ {
			f := st.Field(i)
			fv := rv.Field(i)
			ft := f.Type
			switch ft.Kind() {
			case reflect.Int:
				validationFailedInt(f.Name, int(fv.Int()), f.Tag)
			case reflect.String:
				validationFailedString(f.Name, fv.String(), f.Tag)
			}
		}

	}
	return nil
}
