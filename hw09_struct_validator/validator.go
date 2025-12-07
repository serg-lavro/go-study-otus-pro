package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrMinValidation    = errors.New("min validation failed")
	ErrMaxValidation    = errors.New("max validation failed")
	ErrLenValidation    = errors.New("len validation failed")
	ErrInValidation     = errors.New("in validation failed")
	ErrRegexpValidation = errors.New("regexp validation failed")
	ErrSOME             = errors.New("SOME")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, ve := range v {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(ve.Err.Error())
	}
	return sb.String()
}

type ValidationRequirenment struct {
	reqType string
	reqVal  string
}

func parseRequirenments(tag reflect.StructTag) []ValidationRequirenment {
	validateTag := tag.Get("validate")
	if validateTag == "" {
		return nil
	}
	parts := strings.Split(validateTag, "|")

	reqs := make([]ValidationRequirenment, 0)

	for _, p := range parts {
		reqFields := strings.Split(p, ":")
		reqs = append(reqs, ValidationRequirenment{reqFields[0], reqFields[1]})
	}
	return reqs
}

func validateIn(fieldVal, set string) bool {
	allowedSet := strings.Split(set, ",")
	return slices.Contains(allowedSet, fieldVal)
}

func validationFailedInt(fieldName string, val int, tag reflect.StructTag) (ValidationErrors, bool) {
	errs := ValidationErrors{}
	failed := false
	reqs := parseRequirenments(tag)
	for _, r := range reqs {
		switch r.reqType {
		case "min":
			rval, _ := strconv.Atoi(r.reqVal)
			if rval > val {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrMinValidation,
				})
				failed = true
			}
		case "max":
			rval, _ := strconv.Atoi(r.reqVal)
			if val > rval {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrMaxValidation,
				})
				failed = true
			}
		case "in":
			if !validateIn(strconv.Itoa(val), r.reqVal) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrInValidation,
				})
				failed = true
			}
		}
	}
	return errs, failed
}

func validationFailedString(fieldName string, val string, tag reflect.StructTag) (ValidationErrors, bool) {
	errs := ValidationErrors{}
	failed := false
	reqs := parseRequirenments(tag)
	for _, r := range reqs {
		switch r.reqType {
		case "in":
			if !validateIn(val, r.reqVal) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrInValidation,
				})
				failed = true
			}
		case "len":
			rval, _ := strconv.Atoi(r.reqVal)
			if rval != len(val) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrLenValidation,
				})
				failed = true
			}
		case "regexp":
			matched, _ := regexp.MatchString(r.reqVal, val)
			if !matched {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrRegexpValidation,
				})
				failed = true
			}
		}
	}
	return errs, failed
}

func validateSlice(f reflect.StructField, fv reflect.Value) ValidationErrors {
	var errs ValidationErrors
	elemKind := fv.Type().Elem().Kind()

	switch elemKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < fv.Len(); i++ {
			elem := fv.Index(i)
			ve, failed := validationFailedInt(fmt.Sprintf("%s[%d]", f.Name, i), int(elem.Int()), f.Tag)
			if failed {
				errs = append(errs, ve...)
			}
		}
	case reflect.String:
		for i := 0; i < fv.Len(); i++ {
			elem := fv.Index(i)
			ve, failed := validationFailedString(fmt.Sprintf("%s[%d]", f.Name, i), elem.String(), f.Tag)
			if failed {
				errs = append(errs, ve...)
			}
		}
	case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct,
		reflect.UnsafePointer:
	}

	return errs
}

func Validate(v interface{}) error {
	var errors ValidationErrors
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		st := reflect.TypeOf(v)
		for i := 0; i < st.NumField(); i++ {
			f := st.Field(i)
			fv := rv.Field(i)
			ft := f.Type
			switch ft.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				errs, failed := validationFailedInt(f.Name, int(fv.Int()), f.Tag)
				if failed {
					errors = append(errors, errs...)
				}
			case reflect.String:
				errs, failed := validationFailedString(f.Name, fv.String(), f.Tag)
				if failed {
					errors = append(errors, errs...)
				}
			case reflect.Slice:
				sliceErrs := validateSlice(f, fv)
				errors = append(errors, sliceErrs...)
			case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
				reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func,
				reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct,
				reflect.UnsafePointer:
			default:
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
