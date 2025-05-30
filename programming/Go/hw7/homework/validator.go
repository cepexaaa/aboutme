package homework

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrNotStruct                   = errors.New("wrong argument given, should be a struct")
	ErrInvalidValidatorSyntax      = errors.New("invalid validator syntax")
	ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
	ErrLenValidationFailed         = errors.New("len validation failed")
	ErrInValidationFailed          = errors.New("in validation failed")
	ErrMaxValidationFailed         = errors.New("max validation failed")
	ErrMinValidationFailed         = errors.New("min validation failed")
)

type ValidationError struct {
	field string
	err   error
}

func NewValidationError(err error, field string) error {
	return &ValidationError{
		field: field,
		err:   err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.err)
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

func Validate(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var validationErrors []error
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !field.IsExported() {
			if tag := field.Tag.Get("validate"); tag != "" {
				validationErrors = append(validationErrors,
					NewValidationError(ErrValidateForUnexportedFields, field.Name))
			}
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		var err error
		if v := fieldValue.Interface(); isValidType(v) {
			err = validateTags(v, tag)
		} else if fieldValue.Kind() == reflect.Slice {
			err = validateSlice(fieldValue, tag)
		} else {
			validationErrors = append(validationErrors,
				NewValidationError(ErrValidateForUnexportedFields, field.Name))
		}

		if err != nil {
			validationErrors = append(validationErrors,
				NewValidationError(err, field.Name))
		}
	}

	if len(validationErrors) != 0 {
		return errors.Join(validationErrors...)
	}
	return nil
}

func isValidType[T any](v T) bool {
	switch any(v).(type) {
	case string, int:
		return true
	default:
		return false
	}
}

func isValidSlice(v reflect.Value) bool {
	if v.Kind() != reflect.Slice {
		return false
	}

	elemType := v.Type().Elem()
	switch elemType.Kind() {
	case reflect.String, reflect.Int:
		return true
	default:
		return false
	}
}

func validateSlice(slice reflect.Value, tag string) error {
	var elemErrors []error

	if !isValidSlice(slice) {
		elemErrors = append(elemErrors, errors.New("unsupported type of field"))
		return errors.Join(elemErrors...)
	}
	for i := range slice.Len() {
		elem := slice.Index(i)
		err := validateTags(elem.Interface(), tag)
		if err != nil {
			elemErrors = append(elemErrors, fmt.Errorf("[%d]: %w", i, err))
		}
	}

	if len(elemErrors) > 0 {
		return errors.Join(elemErrors...)
	}
	return nil
}

func validateTags(value any, tag string) error {
	rule, params, err := parseValidator(tag)
	if err != nil {
		return err
	}
	return applyValidationRule(value, rule, params)
}

func parseValidator(validator string) (rule, params string, err error) {
	parts := strings.SplitN(validator, ":", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", "", ErrInvalidValidatorSyntax
	}
	return parts[0], parts[1], nil
}

func applyValidationRule(value any, rule, params string) error {
	switch rule {
	case "len":
		return validateLen(value, params)
	case "in":
		return validateIn(value, params)
	case "min":
		return validateMin(value, params)
	case "max":
		return validateMax(value, params)
	default:
		return ErrInvalidValidatorSyntax
	}
}

func validateMin(value any, params string) error {
	min, err := strconv.Atoi(params)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if getSize(value) < min {
		return ErrMinValidationFailed
	}
	return nil
}

func validateMax(value any, params string) error {
	max, err := strconv.Atoi(params)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if getSize(value) > max {
		return ErrMaxValidationFailed
	}
	return nil
}

func validateIn(value any, params string) error {
	switch v := any(value).(type) {
	case string:
		return checkInToString(v, params)
	case int:
		return checkInToInt(v, params)
	}
	return nil
}

func validateLen(value any, params string) error {
	s, ok := any(value).(string)
	if !ok {
		return errors.New("unsupported parameter 'len' for int")
	}
	length, err := strconv.Atoi(params)
	if err != nil || length < 0 {
		return ErrInvalidValidatorSyntax
	}
	if len(s) != length {
		return ErrLenValidationFailed
	}
	return nil
}

func getSize[T any](value T) int {
	switch v := any(value).(type) {
	case string:
		return len(v)
	case int:
		return v
	default:
		return 0
	}
}

func checkInToInt(value int, params string) error {
	options := strings.Split(params, ",")
	found := false
	for _, opt := range options {
		num, err := strconv.Atoi(opt)
		if err != nil {
			return ErrInValidationFailed
		}
		if num == value {
			found = true
			break
		}
	}
	if !found {
		return ErrInValidationFailed
	}
	return nil
}

func checkInToString(value, params string) error {
	options := strings.Split(params, ",")
	if !slices.Contains(options, value) {
		return ErrInValidationFailed
	}
	return nil
}
