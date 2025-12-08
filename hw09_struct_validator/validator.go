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
	ErrInvalidTag       = errors.New("invalid Tag")
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

func validateInt(fieldName string, val int, tag reflect.StructTag) (ValidationErrors, error) {
	errs := ValidationErrors{}
	reqs := parseRequirenments(tag)
	for _, r := range reqs {
		switch r.reqType {
		case "min":
			rval, e := strconv.Atoi(r.reqVal)
			if e != nil {
				return nil, ErrInvalidTag
			}
			if rval > val {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrMinValidation,
				})
			}
		case "max":
			rval, e := strconv.Atoi(r.reqVal)
			if e != nil {
				return nil, ErrInvalidTag
			}
			if val > rval {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrMaxValidation,
				})
			}
		case "in":
			if !validateIn(strconv.Itoa(val), r.reqVal) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrInValidation,
				})
			}
		default:
			return nil, ErrInvalidTag
		}
	}
	return errs, nil
}

func validateString(fieldName string, val string, tag reflect.StructTag) (ValidationErrors, error) {
	errs := ValidationErrors{}
	reqs := parseRequirenments(tag)
	for _, r := range reqs {
		switch r.reqType {
		case "in":
			if !validateIn(val, r.reqVal) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrInValidation,
				})
			}
		case "len":
			rval, e := strconv.Atoi(r.reqVal)
			if e != nil {
				return nil, ErrInvalidTag
			}
			if rval != len(val) {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrLenValidation,
				})
			}
		case "regexp":
			matched, e := regexp.MatchString(r.reqVal, val)
			if e != nil {
				return nil, ErrInvalidTag
			}
			if !matched {
				errs = append(errs, ValidationError{
					Field: fieldName,
					Err:   ErrRegexpValidation,
				})
			}
		default:
			return nil, ErrInvalidTag
		}
	}
	return errs, nil
}

func validateSlice(f reflect.StructField, fv reflect.Value) (ValidationErrors, error) {
	var errs ValidationErrors
	elemKind := fv.Type().Elem().Kind()

	switch elemKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < fv.Len(); i++ {
			elem := fv.Index(i)
			ve, e := validateInt(fmt.Sprintf("%s[%d]", f.Name, i), int(elem.Int()), f.Tag)
			if e != nil {
				return nil, e
			}
			if ve != nil {
				errs = append(errs, ve...)
			}
		}
	case reflect.String:
		for i := 0; i < fv.Len(); i++ {
			elem := fv.Index(i)
			ve, e := validateString(fmt.Sprintf("%s[%d]", f.Name, i), elem.String(), f.Tag)
			if e != nil {
				return nil, e
			}
			if ve != nil {
				errs = append(errs, ve...)
			}
		}
	case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct,
		reflect.UnsafePointer:
	}

	return errs, nil
}

func Validate(v interface{}) error {
	var errors ValidationErrors
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() != reflect.Struct {
		return nil
	}

	st := reflect.TypeOf(v)
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		fv := rv.Field(i)
		ft := f.Type
		switch ft.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intErrors, e := validateInt(f.Name, int(fv.Int()), f.Tag)
			if e != nil {
				return e
			}
			if intErrors != nil {
				errors = append(errors, intErrors...)
			}
		case reflect.String:
			stringErrors, e := validateString(f.Name, fv.String(), f.Tag)
			if e != nil {
				return e
			}
			if stringErrors != nil {
				errors = append(errors, stringErrors...)
			}
		case reflect.Slice:
			sliceErrs, e := validateSlice(f, fv)
			if e != nil {
				return e
			}
			errors = append(errors, sliceErrs...)
		case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct,
			reflect.UnsafePointer:
		default:
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
