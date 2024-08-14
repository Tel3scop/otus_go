package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for i, ve := range v {
		errs[i] = fmt.Sprintf("%s: %s", ve.Field, ve.Err.Error())
	}
	return strings.Join(errs, "; ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		validators, err := parseValidators(validateTag)
		if err != nil {
			return err
		}

		for _, validator := range validators {
			if err := validateField(field.Name, fieldValue, validator); err != nil {
				var validationErr ValidationError
				if errors.As(err, &validationErr) {
					validationErrors = append(validationErrors, validationErr)
				} else {
					return err
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

type Validator struct {
	Name string
	Arg  string
}

func parseValidators(tag string) ([]Validator, error) {
	validatorStrings := strings.Split(tag, "|")
	validators := make([]Validator, len(validatorStrings))

	for i, validatorString := range validatorStrings {
		parts := strings.SplitN(validatorString, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid validator format: %s", validatorString)
		}
		validators[i] = Validator{Name: parts[0], Arg: parts[1]}
	}

	return validators, nil
}

func validateString(fieldName, fieldValue, validatorName, validatorArg string) error {
	switch validatorName {
	case "len":
		length, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidLength, validatorArg)
		}
		if len(fieldValue) != length {
			return ValidationError{Field: fieldName, Err: fmt.Errorf("length must be %d", length)}
		}
	case "minLen":
		minLength, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidMinLen, validatorArg)
		}
		if len(fieldValue) < minLength {
			return ValidationError{Field: fieldName, Err: fmt.Errorf("minimum length is %d", minLength)}
		}
	case "regexp":
		re, err := regexp.Compile(validatorArg)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidRegexp, validatorArg)
		}
		if !re.MatchString(fieldValue) {
			return ValidationError{Field: fieldName, Err: fmt.Errorf("must match regexp %s", validatorArg)}
		}
	case "in":
		allowedValues := strings.Split(validatorArg, ",")
		for _, val := range allowedValues {
			if fieldValue == val {
				return nil
			}
		}
		return ValidationError{Field: fieldName, Err: fmt.Errorf("must be one of %s", validatorArg)}
	default:
		return fmt.Errorf("%w: %s", errUnknownValidator, validatorName)
	}
	return nil
}

func validateInt(fieldName string, fieldValue int, validatorName, validatorArg string) error {
	switch validatorName {
	case "min":
		minValue, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidMinValue, validatorArg)
		}
		if fieldValue < minValue {
			return ValidationError{Field: fieldName, Err: fmt.Errorf("must be at least %d", minValue)}
		}
	case "max":
		maxValue, err := strconv.Atoi(validatorArg)
		if err != nil {
			return fmt.Errorf("%w: %s", errInvalidMaxValue, validatorArg)
		}
		if fieldValue > maxValue {
			return ValidationError{Field: fieldName, Err: fmt.Errorf("must be at most %d", maxValue)}
		}
	case "in":
		allowedValues := strings.Split(validatorArg, ",")
		for _, val := range allowedValues {
			if intVal, err := strconv.Atoi(val); err == nil && fieldValue == intVal {
				return nil
			}
		}
		return ValidationError{Field: fieldName, Err: fmt.Errorf("must be one of %s", validatorArg)}
	default:
		return fmt.Errorf("%w: %s", errUnknownValidator, validatorName)
	}
	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, validator Validator) error {
	switch fieldValue.Kind() {
	case reflect.String:
		return validateString(fieldName, fieldValue.String(), validator.Name, validator.Arg)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateInt(fieldName, int(fieldValue.Int()), validator.Name, validator.Arg)
	case reflect.Slice:
		if fieldValue.Type().Elem().Kind() == reflect.Uint8 {
			return validateString(fieldName, string(fieldValue.Bytes()), validator.Name, validator.Arg)
		}
		return validateSlice(fieldName, fieldValue, validator.Name, validator.Arg)
	case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
		return fmt.Errorf("%w: %s", errUnsupportedType, fieldValue.Kind())
	default:
		return fmt.Errorf("%w: %s", errUnsupportedType, fieldValue.Kind())
	}
}

func validateSlice(fieldName string, fieldValue reflect.Value, validatorName, validatorArg string) error {
	for i := 0; i < fieldValue.Len(); i++ {
		elem := fieldValue.Index(i)
		switch elem.Kind() {
		case reflect.String:
			if err := validateString(fieldName, elem.String(), validatorName, validatorArg); err != nil {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if err := validateInt(fieldName, int(elem.Int()), validatorName, validatorArg); err != nil {
				return err
			}
		case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array,
			reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct,
			reflect.UnsafePointer:
			return fmt.Errorf("unsupported element type in slice: %s", elem.Kind())
		default:
			return fmt.Errorf("unknown element type in slice: %s", elem.Kind())
		}
	}
	return nil
}
